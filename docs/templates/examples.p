<html>
    <head>
        <meta charset="utf-8">
        <title>Examples</title>
    </head>
    <body>
        <span>Name from cookie: {{.Name}}</span>
        <form action="/upload" method="post">
            Select image to upload:
            <input type="text" name="name">
            <input type="submit" value="Upload Name" name="submit">
        </form>
    </body>
</html>