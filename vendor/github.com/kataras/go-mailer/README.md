<a href="https://travis-ci.org/kataras/go-mailer"><img src="https://img.shields.io/travis/kataras/go-mailer.svg?style=flat-square" alt="Build Status"></a>
<a href="https://github.com/kataras/go-mailer.v0/blob/master/LICENSE"><img src="https://img.shields.io/badge/%20license-MIT%20%20License%20-E91E63.svg?style=flat-square" alt="License"></a>
<a href="https://github.com/kataras/go-mailer.v0/releases"><img src="https://img.shields.io/badge/%20release%20-%20v0.0.2-blue.svg?style=flat-square" alt="Releases"></a>
<a href="#docs"><img src="https://img.shields.io/badge/%20docs-reference-5272B4.svg?style=flat-square" alt="Read me docs"></a>
<a href="https://kataras.rocket.chat/channel/go-mailer"><img src="https://img.shields.io/badge/%20community-chat-00BCD4.svg?style=flat-square" alt="Build Status"></a>
<a href="https://golang.org"><img src="https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square" alt="Built with GoLang"></a>
<a href="#"><img src="https://img.shields.io/badge/platform-Any--OS-yellow.svg?style=flat-square" alt="Platforms"></a>


Simple E-mail sender for Go Programming Language.

Supports rich e-mails and ,optionally, Unix built'n  `sendmail` command.

Installation
------------
The only requirement is the [Go Programming Language](https://golang.org/dl).

```bash
$ go get -u gopkg.in/kataras/go-mailer.v0
```


Docs
------------

- `New` returns a new, e-mail sender service.
- `Send` send an e-mail, supports text/html and `sendmail` unix command
```go
Send(subject string, body string, to ...string) error
```

**Configuration**

```go
type Config struct {
    // Host is the server mail host, IP or address
    Host string
    // Port is the listening port
    Port int
    // Username is the auth username@domain.com for the sender
    Username string
    // Password is the auth password for the sender
    Password string
    // FromAddr is the 'from' part of the mail header, it overrides the username
    FromAddr string
    // FromAlias is the from part, if empty this is the first part before @ from the Username field
    FromAlias string
    // UseCommand enable it if you want to send e-mail with the mail command  instead of smtp
    //
    // Host,Port & Password will be ignored
    // ONLY FOR UNIX
    UseCommand bool
}

```

**Usage**

```go
import "gopkg.in/kataras/go-mailer.v0"

//...

// sender configuration
config := mailer.Config{
    Host:     "smtp.mailgun.org",
    Username: "postmaster",
    Password: "38304272b8ee5c176d5961dc155b2417",
    FromAddr: "postmaster@sandbox661c307650f04e909150b37c0f3b2f09.mailgun.org",
    Port:     587,
    UseCommand:false,
    // Enable UseCommand to support sendmail unix command, if this field is true then Host, Username, Password and Port are not required, because these info already exists in your local sendmail configuration
}

// initalize a new mail sender service
sender := mailer.New(config)

// the subject/title of the e-mail
subject := "Hello subject"

// the rich message body
content := `<h1>Hello</h1> <br/><br/> <span style="color:red"> This is the rich message body </span>`

// the recipient(s)
to := []string{"kataras2006@hotmail.com", "social@ideopod.com"}


// Send the e-mail
err := mailer.Send(subject, content, to...)

if err != nil {
  println("Error while sending the e-mail. Trace: "+err.Error())
}

```

FAQ
------------
Explore [these questions](https://gopkg.in/kataras/go-mailer.v0/issues?go-mailer=label%3Aquestion) or navigate to the [community chat][Chat].

Versioning
------------

Current: **v0.0.2**



People
------------
The author of go-mailer is [@kataras](https://github.com/kataras).

If you're **willing to donate**, feel free to send **any** amount through paypal

[![](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/cgi-bin/webscr?cmd=_donations&business=kataras2006%40hotmail%2ecom&lc=GR&item_name=Iris%20web%20framework&item_number=iriswebframeworkdonationid2016&currency_code=EUR&bn=PP%2dDonationsBF%3abtn_donateCC_LG%2egif%3aNonHosted)


Contributing
------------
If you are interested in contributing to the go-mailer project, please make a PR.

License
------------

This project is licensed under the MIT License.

License can be found [here](LICENSE).

[Travis Widget]: https://img.shields.io/travis/kataras/go-mailer.svg?style=flat-square
[Travis]: http://travis-ci.org/kataras/go-mailer
[License Widget]: https://img.shields.io/badge/license-MIT%20%20License%20-E91E63.svg?style=flat-square
[License]: https://gopkg.in/kataras/go-mailer.v0/blob/master/LICENSE
[Release Widget]: https://img.shields.io/badge/release-v4.1.1-blue.svg?style=flat-square
[Release]: https://gopkg.in/kataras/go-mailer.v0/releases
[Chat Widget]: https://img.shields.io/badge/community-chat-00BCD4.svg?style=flat-square
[Chat]: https://kataras.rocket.chat/channel/go-mailer
[ChatMain]: https://kataras.rocket.chat/channel/go-mailer
[ChatAlternative]: https://gitter.im/kataras/go-mailer
[Report Widget]: https://img.shields.io/badge/report%20card-A%2B-F44336.svg?style=flat-square
[Report]: http://goreportcard.com/report/kataras/go-mailer
[Documentation Widget]: https://img.shields.io/badge/documentation-reference-5272B4.svg?style=flat-square
[Documentation]: https://www.gitbook.com/book/kataras/go-mailer/details
[Language Widget]: https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square
[Language]: http://golang.org
[Platform Widget]: https://img.shields.io/badge/platform-Any--OS-gray.svg?style=flat-square
