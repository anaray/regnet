regnet 
======
Status : Work In Progress


Complex regular expressions from simple manageable regex units

Example:

```
DAY = `(?:Mon(?:day)?|Tue(?:sday)?|Wed(?:nesday)?|Thu(?:rsday)?|Fri(?:day)?|Sat(?:urday)?|Sun(?:day)?)`
```

```
MONTH = `(?:Jan(?:uary)?|Feb(?:ruary)?|Mar(?:ch)?|Apr(?:il)?|May|Jun(?:e)?|Jul(?:y)?|Aug(?:ust)?|Sep(?:tember)?|Oct(?:ober)?|Nov(?:ember)?|Dec(?:ember)?)`
```

```
YEAR = `(\d\d){1,2}`
```
So by combining the regnets
%{DAY} %{MONTH} %{YEAR}

can match a text **"This note was written on Tue Aug 2014"** to give **"Tue Aug 2014"**

Usage:
```
  r, _ := regnet.New()
  r.AddPattern("DAY", `(?:Mon(?:day)?|Tue(?:sday)?|Wed(?:nesday)?|Thu(?:rsday)?|Fri(?:day)?|Sat(?:urday)?|Sun(?:day)?)`)
  r.AddPattern("YEAR", `(\d\d){1,2}`)
  r.AddPattern("ACTION_DATE", `%{DAY} May %{YEAR}`)
  
  //regnets can be defined in a file and can be loaded as below
  //err := r.AddPatternsFromFile("/regex/text/patterns/extra")

  match, _ := r.MatchInText("Tue May 15 11:21:42 [conn1047685] moveChunk deleted: 7157", "%{ACTION_DATE}")

  for _, result := range match.Step() {
    fmt.Println(match.Ident, ":" ,result)
  }
	
```
Inspired by Jordan Sissel's work https://github.com/jordansissel/grok