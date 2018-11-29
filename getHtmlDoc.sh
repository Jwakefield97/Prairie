#!/bin/sh
#make sure godocs server is running. Use command (godoc -html :6060)
mkdir docs
mkdir docs/css
mkdir docs/js

#html
godoc -url http://localhost:6060/pkg/github.com/Jwakefield97/prairie/ > docs/prairie.html
godoc -url http://localhost:6060/pkg/github.com/Jwakefield97/prairie/lib/utils/ > docs/utils.html
godoc -url http://localhost:6060/pkg/github.com/Jwakefield97/prairie/lib/http/ > docs/http.html

#css
godoc -url http://localhost:6060/lib/godoc/style.css > docs/css/style.css
godoc -url http://localhost:6060/lib/godoc/jquery.treeview.css > docs/css/jquery.treeview.css

#js
godoc -url http://localhost:6060/lib/godoc/jquery.js > docs/js/jquery.js
godoc -url http://localhost:6060/lib/godoc/jquery.treeview.js > docs/js/jquery.treeview.js
godoc -url http://localhost:6060/lib/godoc/jquery.treeview.edit.js > docs/js/jquery.treeview.edit.js
godoc -url http://localhost:6060/lib/godoc/godocs.js > docs/js/godocs.js