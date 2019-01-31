<html>
    <head>
        <meta charset="utf-8">
        <title>Examples</title>
    </head>
    <body>
        <span>Name from cookie: {{.Name}}</span>
        <form action="/upload" method="post">
            Type name to upload:
            <input type="text" name="name">
            <input type="submit" value="Upload Name" name="submit">
        </form>

        <code>
            <a href="/template">render the below template</a> 
            <pre>
                app.Get("/template", func(routeObj *prairie.RouteObject) {
                    routeObj.Response.Template = "temp"
                    routeObj.Response.TemplateParams = TodoPageData{
                        PageTitle: "My TODO list",
                        Todos: []Todo{
                            {Title: "Task 1", Done: false},
                            {Title: "Task 2", Done: true},
                            {Title: "Task 3", Done: true},
                        },
                    }
                })
            </pre>
        </code>
    </body>
</html>