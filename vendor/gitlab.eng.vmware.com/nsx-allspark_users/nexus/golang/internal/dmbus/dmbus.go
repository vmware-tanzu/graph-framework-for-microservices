package dmbus

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/internal/graphdb"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/internal/utils"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/common"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/ifc"
	"gitlab.eng.vmware.com/nsx-allspark_users/nexus/golang/pkg/logging"

	"github.com/lithammer/shortuuid"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

type DMBus struct {
	stats              ifc.NotificationStats
	cbfn               *func(msg *ifc.Notification)
	msgDelay           uint32
	running            bool
	createRotationTime int
	createMaxSize      int
	createDataPtr      bool
	createDataA        map[string]string
	createDataB        map[string]string
	name               string
	dmId               string
}

func New(name string, dmId string) *DMBus {
	e := &DMBus{
		createRotationTime: 60 * 60 * 1000, // 1hour
		createMaxSize:      100000,         // 100k entries max
		createDataPtr:      false,
		msgDelay:           0,
		name:               name,
		dmId:               dmId,
		stats:              ifc.NotificationStats{},
		createDataA:        make(map[string]string),
		createDataB:        make(map[string]string),
		running:            false,
	}
	logging.Debugf("CREATE: new DMBUS %s[%+v] DONE\n", name, *e)
	return e
}
func (d *DMBus) Init(graph interface{}) {
	g := graph.(*graphdb.GraphDB)
	logging.Debugf("%s:Starting watch on Root", d.name)
	rch := g.GetClient().Watch(context.Background(), "/Root",
		clientv3.WithPrefix())
	d.running = true // This field needs to be made atomic. Potential race here.
	// store rch for shutdown implementation
	go func() {
		cnt := 0
		logging.Debugf("DataModel[%s] %p starting watch for events\n", d.name, d)
		for itm := range rch {
			if d.running {
				cnt++
				d.data(itm)
			} else {
				break
			}
		}
		logging.Debugf("SHUTDOWN:DataModel[%s] %p exiting watch for events %d\n", d.name, d, cnt)
	}()
}
func (d *DMBus) data(wresp clientv3.WatchResponse) {
	unknownWarnMap := make(map[string]struct{})
	for _, msg := range wresp.Events {
		key := string(msg.Kv.Key)
		valueStr := string(msg.Kv.Value)
		var mtype mvccpb.Event_EventType = msg.Type
		pk := NewProcessKey().process(key)

		logging.Debugf("%s:Got notification[%d] type=%s  key = %s ==> %s",
			d.name, msg.Kv.ModRevision, mtype, key, valueStr)

		if pk.UnknownNodeType {
			if _, ok := unknownWarnMap[key]; !ok {
				/*Should be set to warning log level, but we'd end up observing a flood of such warnings.
				  Looks like we end up receving multiple notifications for each update into the database,
				  in general. Needs to be evaluated.
				*/
				logging.Debugf("%s:Will Skip as unknown node type is identified for key %s",
					d.name, key)
				/* Need this check to prevent flooding of warning messages */
				unknownWarnMap[key] = struct{}{}
			}
			continue
		}

		if pk.Created || pk.Lock {
			logging.Debugf("%s:Will Skip as it is created/lock", d.name)
			if pk.Created {
				d.addCreateData(key, valueStr)
			}
			continue
		}
		nf := &ifc.Notification{
			ObjectPath:         pk.Path,
			Source:             d.dmId,
			Timestamp:          time.Now().UTC().Format(utils.ISOTimeFormat),
			UpdateType:         "",
			UpdatedObjId:       pk.NodeId,
			UpdatedObjType:     pk.NodeType,
			UpdatedObjKey:      pk.NodeKey,
			UpdatedObjParentId: pk.ParentId,
			TraceId:            fmt.Sprintf("<M:%s:%s>", shortuuid.New(), pk.Path),
			Value:              make(ifc.PropertyType),
		}
		nf.Revision = msg.Kv.ModRevision
		if valueStr != "" {
			if err := json.Unmarshal(msg.Kv.Value, &nf.Value); err != nil {
				logging.Fatalf("Received Malformed JSON in etcd watch: %s %s", valueStr, err)
			}
		}
		if pk.Link {
			nf.UpdatedObjKey = pk.LinkInfoDestNodeKeyValue
			nf.UpdatedObjType = pk.LinkInfoDestNodeType
			nf.UpdatedObjParentId = pk.NodeId
			nf.UpdatedObjId = key // pk.nodeId;
			nf.UpdateType = common.UpdateType_StrLinkUpdate
			if mtype == mvccpb.DELETE {
				nf.UpdateType = common.UpdateType_StrLinkDelete
			}
			cpath := key + "/" + graphdb.FixedPath_Created
			cdata := d.checkCreateData(cpath)
			// if create node is not received yet
			// if cdata == "" {
			// 	time.Sleep(time.Duration(100) * time.Millisecond)
			// 	cdata = d.checkCreateData(cpath)
			// }
			if cdata != "" && nf.Value != nil {
				var cd ifc.PropertyType
				if err := json.Unmarshal([]byte(cdata), &cd); err != nil {
					logging.Debugf("Error parsing cdata %s", err)
				}
				nf.Value[common.LinkFixedProp_createdBy] = cd[common.LinkFixedProp_createdBy]
				nf.Value[common.LinkFixedProp_creationTime] = cd[common.LinkFixedProp_creationTime]
			}
			d.stats.LinkUpdateRx++
		} else if pk.RLink {
			/*
				IMPORTANT PLEASE READ:
				We dont want to generate an event at a node level here or at a link level above
				to preserve backward compatibility.
				Add events for rlinks would impact consumers subscribed to node events.
			*/
			continue
		} else {
			nf.UpdateType = common.UpdateType_StrNodeUpdate
			if mtype == mvccpb.DELETE {
				nf.UpdateType = common.UpdateType_StrNodeDelete
			}
			cpath := key + "/" + graphdb.FixedPath_Created
			cdata := d.checkCreateData(cpath)
			if cdata != "" && nf.Value != nil {
				var cd ifc.PropertyType
				if err := json.Unmarshal([]byte(cdata), &cd); err != nil {
					logging.Debugf("Error parsing cdata %s", err)
				}
				nf.Value[common.NodeFixedProp_createdBy] = cd[common.NodeFixedProp_createdBy]
				nf.Value[common.NodeFixedProp_creationTime] = cd[common.NodeFixedProp_creationTime]
			}
			d.stats.NodeUpdateRx++
		}

		// TODO: need a scheduler that lets us synchronize on node/link ID of events coming in.
		// We want to allow for parallelism but not for the same node at the same time.
		go func(nf *ifc.Notification) {
			if d.msgDelay != 0 {
				time.Sleep(time.Duration(d.msgDelay) * time.Millisecond)
			}

			logging.Debugf("%s:Notification: %v", d.name, nf)
			if d.cbfn != nil {
				(*d.cbfn)(nf)
			}
		}(nf)
	}
}
func (d *DMBus) RegisterCB(fn func(msg *ifc.Notification)) {
	d.cbfn = &fn
}
func (d *DMBus) AddMessageDelay(n uint32) {
	logging.Warnf("%s:MessageBusETC Setup to create addition message delivery delay of %d ms", d.name, n)
	d.msgDelay = n
}
func (d *DMBus) toggleDataPtr() {
	if d.createDataPtr {
		d.createDataA = make(map[string]string)
		d.createDataPtr = false
	} else {
		d.createDataB = make(map[string]string)
		d.createDataPtr = true
	}

	// FIXME: refactor the below. This is infinitely recursive, will eventually stack overflow after days/weeks/months of runtime?
	//d.createDataPtrTimeout = setTimeout(this.toggleDataPtr.bind(this), this.createRotationTime);
	fn := func() {
		time.Sleep(time.Duration(d.createRotationTime) * time.Millisecond)
		d.toggleDataPtr()
	}
	if d.running {
		go fn()
	}
}
func (d *DMBus) addCreateData(key, value string) {
	if d.createDataPtr {
		d.createDataB[key] = value
		if len(d.createDataB) > d.createMaxSize {
			// clearTimeout(d.createDataPtrTimeout);
			d.toggleDataPtr()
		}
	} else {
		d.createDataA[key] = value
		if len(d.createDataA) > d.createMaxSize {
			// clearTimeout(d.createDataPtrTimeout);
			d.toggleDataPtr()
		}
	}
}
func (d *DMBus) checkCreateData(key string) string {
	r, ok := d.createDataA[key]
	if !ok {
		r, ok = d.createDataB[key]
	}
	if ok {
		return r
	} else {
		return ""
	}
}
func (d *DMBus) Shutdown() {
	logging.Infof("SHUTDOWN: Stopping DMBUS %s[%p]\n", d.name, d)
	d.running = false
}
func (d *DMBus) GetStats() ifc.NotificationStats {
	return d.stats
}
