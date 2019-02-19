File Component
==============

Functionality:
--------------

A domain relates to a directory within a filesystem.

Each component has a name which is unique within a domain.
It resides at a location relative to its parent component within a file system and holds a template which may contain references to other components, eg.


<code>
Hello <span style="color:orange">{{customer}}</span>&nbsp;   
</code>

The reference {{customer}} in this case would relate to the name of another component which resides in another location within the same filesystem.

The functionality is provided by a simple web service which is given a path in the filesystem and then recursively evaluates the templates as needed and returns the result.

Special note: the evaluation algorithm evaluates all instances of a component and makes use of the majority opinion to simulate cluster behaviour.

General structure
-----------------

Files are represented in a directory structure within a root context:

```
<root>                      Directory
  ...                       ...
    <parent>                Directory    (relates to a runtime context)
      <element>             Directory    (relates to an element)
        .VX.Y.Z             Directory    (relates to a cluster)
          <instance 1>      File         (relates to an instance)
          <instance 2>      File         (relates to an instance)
          <instance 3>      File         (relates to an instance)
          <instance ...>    File         (relates to an instance)
        <child 1>           Directory
        <child 2>           Directory
        <child 3>           Directory
        <child ...>         Directory
      <sibling 1>           Directory
      <sibling 2>           Directory
      <sibling 3>           Directory
      <sibling ...>         Directory
```

Information held in <instance N>
--------------------------------

Content is in yaml format:

```
state:    active/inactive
path:     <path>
template: <template string>
references:
  - <name 1>: <path 1>
  - <name 2>: <path 2>
  - <name 3>: <path 3>
  ...
```
