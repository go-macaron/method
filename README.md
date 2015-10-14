method
=======

Middlware method implements HTTP method override for [Macaron](https://github.com/go-macaron/macaron).

This checks for the X-HTTP-Method-Override header and uses it
if the original request method is POST.
GET/HEAD methods shouldn't be overriden, hence they can't be overriden.

This is useful for REST APIs and services making use of many HTTP verbs, and when http clients don't support all of them.

[API Reference](https://gowalker.org/github.com/go-macaron/method)

## Usage

```go
import (
	"gopkg.in/macaron.v1"
  "github.com/go-macaron/method"
)

func main() {
  m := macaron.Classic()
  m.Before(method.Override())
  m.Run()
}
```

## Credits

This package is forked from [martini-contrib/method](https://github.com/martini-contrib/method) with modifications.

## License

This project is under Apache v2 License. See the [LICENSE](LICENSE) file for the full license text.
