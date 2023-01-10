# Nexus shim layer

Nexus client is a shim layer working over k8s API meant to simplify graph operation on nexus nodes. Features of nexus
client:

- create/get/update/delete/list of nexus nodes,
- name hashing to avoid name collision between objects with same name but different parents,
- ability to get, create and delete child of given parent object,
- ability to add link and remove link to given object,
- recursive delete of object and all it's children.
- supports get/set/clear for user defined status of a nexus node.  

# API

Currently following API calls are available:

- `nexusClient.{Group Type}{Object Type}.Add(context, objToCreate)` creates root type nexusObject (root objects are
  objects which don't have any parents). It calculates hashed name of the object based on objToCreate.Name and creates
  it. objToCreate.Name is changed to the hashed name. Original name is preserved in nexus/display_name label and can be
  obtained using DisplayName() method.
- `nexusClient.{Group Type}{Object Type}.Delete(context, displayName)` deletes root type nexusObject.
- `nexusClient.{Group Type}{Object Type}.Get(context, displayName)` returns given root type nexusObject.
- `{nexusObject}.Delete(context)` is a method on `nexusObject`, removes `nexusObject` and all it's children from the
  database.
- `{nexusObject}.Update(context)` is a method on `nexusObject`, updates spec of `nexusObject` in database. Children and
  Links can not be updated using this function.
- `{nexusObject}.Get{Child}(context)` is a method on `nexusObject`, it returns given child.
- `{nexusObject}.Get{Child}(context, displayName)` is a method on `nexusObject` when child is named. It returns child
  with given name.
- `{nexusObject}.GetAll{Child}(context, displayName)` is a method on `nexusObject` when child is named. It returns all
  children of this type.
- `{nexusObject}.Add{Child}(context, objToCreate)` is a method on `nexusObject` which creates child objToCreate. It
  calculates hashed name of the child to create based on objToCreate.Name and parents names and creates it.
  objToCreate.Name is changed to the hashed name. Original name is preserved in nexus/display_name label and can be
  obtained using DisplayName() method.
- `{nexusObject}.Delete{Child}(context)` is a method on `nexusObject` when this is single child. It removes this child.
- `{nexusObject}.Delete{Child}(context, displayName)` is a method on `nexusObject` when child is named. It removes given
  child.
- `nexusClient.{Group Type}.Get{Object Type}ByName(context, hashedName)` returns nexus object stored in the database
  under the hashedName which is a hash of display name and parents names. Use it when you know hashed name of object.
- `{nexusObject}.Link{Softlink}(context, linkToAdd)` is a method on `nexusObject` which allows adding link to linkToAdd
  object. This function doesn't create linked object, it must be already created.
- `{nexusObject}.Unlink{Softlink}(context)` is a method on `nexusObject` which allows removing link to object when this
  is single object link. This function doesn't delete linked object.
- `{nexusObject}.Unlink{Softlink}(context, linkToRemove)` is a method on `nexusObject` which allows removing link to
  `linkToRemove` object when this is named object link. This function doesn't delete linked object.
- `{nexusObject}.GetParent(context)` is a method on `nexusObject` which allows getting parent object of `nexusObject`.
- `nexusClient.{Group Type}.Delete{Object Type}ByName(context, hashedName)` deletes nexus object stored in the database
  under the hashedName which is a hash of display name and parents names. Use it when you know hashed name of object.
- `nexusClient.{Group Type}.Create{Object Type}ByName(context, objToCreate)` creates nexus object in the database
  without hashing the name. Use it directly ONLY when objToCreate.Name is hashed name of the object.
- `nexusClient.{Group Type}.Update{Object Type}ByName(context, objToUpdate)` updates nexus object stored in the database
  under the hashedName which is a hash of display name and parents names.
- `nexusClient.{Group Type}.List{Object Type}(context, listOptions)` returns slice of all existing nexus objects of this
  type. Selectors can be provided in opts parameter.
- `nexusClient.{Root Object Type}(name).{Child1 of Root Object Type}(name).{Child2 of Child1}(name)...{ChildN of ChildN-1}.Get{Object}(context, displayName)`
  is a method which allows getting nexus object based on provided hierarchy.
- `nexusClient.{Root Object Type}(name).{Child1 of Root Object Type}(name).{Child2 of Child1}(name)...{ChildN of ChildN-1}.Add{Object}(context, displayName)`
  is a method which allows creating child based on provided hierarchy.
- `nexusClient.{Root Object Type}(name).{Child1 of Root Object Type}(name).{Child2 of Child1}(name)...{ChildN of ChildN-1}.Delete{Object}(context, displayName)`
  is a method which allows deleting child based on provided hierarchy.

Available status APIs

- `nexusClient.{Group Type}.Set{StatusType}ByName(context, nexusObjectStatusToSet)` sets nexus object's status stored in the database
  under the hashedName which is a hash of display name and parents names.
- `{nexusObject}.Get{Object StatusName}(context)` is a method on `nexusObject`, to get nexus object's status.
- `{nexusObject}.Set{Object StatusName}(context, nexusObjectStatusToSet)` is a method on `nexusObject`, to set nexus object's status. 
- `{nexusObject}.Clear{Object StatusName}(context)` is a method on `nexusObject`, to clear nexus object's status.
- `nexusClient.{Root Object Type}(name).{Child1 of Root Object Type}(name).{Child2 of Child1}(name)...{ChildN of ChildN-1}(name).Get{Object Status Name}(context)`
  is a method which allows getting nexus object's status based on provided hierarchy.
- `nexusClient.{Root Object Type}(name).{Child1 of Root Object Type}(name).{Child2 of Child1}(name)...{ChildN of ChildN-1}(name).Set{Object Status Name}(context, nexusObjectStatusToSet)`
  is a method which allows setting nexus object's status based on provided hierarchy.
- `nexusClient.{Root Object Type}(name).{Child1 of Root Object Type}(name).{Child2 of Child1}(name)...{ChildN of ChildN-1}(name).Clear{Object Status Name}(context)`
  is a method which allows clearing nexus object's status based on provided hierarchy.


# Working with shim layer

To initialize client use NewForConfig function with Rest Config as a parameter. After that you can start using nexus
client. You can check example app
[here](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/docs/-/tree/master/example/crudapp). Another app
showing how shim layer can be used in operator framework can be found
[here](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/docs/-/tree/master/example/operatorapp).

# Mock client

In unit tests you can use fake client initialized by NewFakeClient function. This client doesn't require database, all
operations are done in memory. You can check example in compiler's
[unit tests](https://gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/compiler/-/blob/master/example/tests/nexusclient_test.go)
.
