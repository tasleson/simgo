### simgo
LibStorageMgmt simulator plugin written in go

This is an experimental plugin for libStorageMgmt using golang library from
https://github.com/libstorage/libstoragemgmt-golang

The plugin is actually just a forwarding plugin in that it take the requests
and forwards them to a different plugin.  This is done by setting the URI to
contain the query string "forward=<plugin>".  An example:


```bash
export LSMCLI_URI='simgo://ignore?forward=sim'
```
