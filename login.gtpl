{{define "Index"}}
<!DOCTYPE html>
<html>
<div class="container">
    <div class="row">
        <div class="col-xs-6 col-xs-offset-4">
            <h3>Sign Up: $10</h3>
            <form method="POST">
                <p>
                    <label for="email">email</label><br>
                    <input type="text" name="email" autofocus="autofocus">
                </p>
                <p>
                    <label for="password">Password</label><br>
                    <input type="password" name="password">
                </p>
                <input type="submit" name="commit" value="Sign Up" class="btn btn-default">
            </form>
        </div>
    </div>
</div>
</html>
{{end}}
