# Conditional Redirect

## Configuration

```yaml
statusCode: 302 # The http status code of redirection response. default: 302
rules: [] # The rules of redirection
```

### Rule

```yaml
# When request matches this pattern, will check this rule for redirection.
# Caveat: it must be a regular expression so you should assert beginning and end correctly
sourcePattern: "^/(.*)"
# When set to `true`, sourcePattern will match full url includes host and protocol
withHost: false
# The target url when condition is met.
# It is able to use groups matched by source pattern to do replacement
destinationPattern: "/foo/$1"
# The condition to check.
condition:
  type: header
  name: foo
  pattern: ".*"
```

### Condition

#### Header

```yaml
type: header
# When set to true, without this header will also pass the check
optional: true
# Header name 
name: foo
# Pattern to match header value
pattern: ".*"
```

#### Cookie

```yaml
type: cookie
# When set to true, without this header will also pass the check
optional: true
# Header name 
name: foo
# The path of the cookie, it is not a regular expression
path: ".*"
# Pattern to match header value
pattern: ".*"
```

#### Logic

```yaml
# Only when all children of it passed, it will pass
type: and
# Conditions
children: []
```

```yaml
# When any child of it passed, it will pass
type: or
# Conditions
children: []
```

```yaml
# Negate the result of condition
type: not
condition:
  type: header
  name: foo
  pattern: ".*"
```

