<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} | TravelSphere</title>
    <meta name="description" content="{{.MetaDescription}}">
    <link rel="stylesheet" href="/static/css/main.css">
</head>
<body>
    {{template "partials/header.tpl" .}}

    <main class="main-content">
        {{.LayoutContent}}
    </main>

    {{template "partials/footer.tpl" .}}

    <script src="/static/js/main.js"></script>
    {{if .PageScript}}
    <script src="/static/js/{{.PageScript}}"></script>
    {{end}}
</body>
</html>