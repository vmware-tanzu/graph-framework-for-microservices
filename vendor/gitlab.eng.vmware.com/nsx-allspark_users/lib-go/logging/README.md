logging usage:

import "gitlab.eng.vmware.com/nsx-allspark_users/lib-go/logging"

        logging.Debugf("Debug log")
        logging.Infof("Info log, my message: %s, number %d", "test msg", 3)
        logging.Errorf("Error log")
        logging.Fatalf("Fatal log")

Output example:

2018-12-12T17:12:47.729-0800	DEBUG	logging/unit_test.go:35	zap Debug log	{"package": "testing", "function": "tRunner"}
2018-12-12T17:12:47.729-0800	INFO	logging/unit_test.go:42	zap Info log msg	{"package": "testing", "function": "tRunner"}
2018-12-12T17:12:47.729-0800	WARN	logging/unit_test.go:48	zap Warn log msg	{"package": "testing", "function": "tRunner"}
2018-12-12T17:12:47.729-0800	ERROR	logging/unit_test.go:54	zap Error log msg	{"package": "testing", "function": "tRunner"}
