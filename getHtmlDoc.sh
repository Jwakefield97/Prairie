#!/bin/sh
#make sure godocs server is running. Use command (godoc -html :6060)

#html
godoc -url http://localhost:6060/pkg/github.com/Jwakefield97/prairie/ > docs/templates/prairie.html
godoc -url http://localhost:6060/pkg/github.com/Jwakefield97/prairie/lib/utils/ > docs/templates/utils.html
godoc -url http://localhost:6060/pkg/github.com/Jwakefield97/prairie/lib/http/ > docs/templates/http.html

#css
godoc -url http://localhost:6060/lib/godoc/style.css > docs/resources/css/style.css
godoc -url http://localhost:6060/lib/godoc/jquery.treeview.css > docs/resources/css/jquery.treeview.css

#js
godoc -url http://localhost:6060/lib/godoc/jquery.js > docs/resources/js/jquery.js
godoc -url http://localhost:6060/lib/godoc/jquery.treeview.js > docs/resources/js/jquery.treeview.js
godoc -url http://localhost:6060/lib/godoc/jquery.treeview.edit.js > docs/resources/js/jquery.treeview.edit.js
godoc -url http://localhost:6060/lib/godoc/godocs.js > docs/resources/js/godocs.js