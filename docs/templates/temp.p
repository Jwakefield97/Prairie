<html>
    <head>
        <meta charset="utf-8">
        <style>
            .done {
                text-decoration: line-through;
            }
        </style>
    </head>
    <body>
        <h1>{{.PageTitle}}<h1>
        <ul>
            {{range .Todos}}
                {{if .Done}}
                    <li class="done">{{.Title}}</li>
                {{else}}
                    <li>{{.Title}}</li>
                {{end}}
            {{end}}
        </ul>
    </body>
</html>


