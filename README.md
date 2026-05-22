# Rapela
Daily prompt to take a moment to pray 🙇🏾

# How It Works
By default once you have installed the tool it will display a bible verse for the day.
For example: `(en-kjv) Genesis 1:1 - In the beginning God created the heaven and the earth.`

Should you clear your terminal then you can retrieve the bible verse of the day by executing the following command:
```bash
rapela verse
```

To see the current bible version:
```bash
rapela bible
```

To see the list of bible versions:
```bash
rapela bible --list-versions or --lv
```

To set a bible version:
```bash
rapela bible --set=<name-of-version>
```

## Future CLI Functionality

To set a language:
```bash
rapela bible --lang=<language-identifier>
```

To log a prayer session:
```bash
rapela prayer
```

To view past prayers:
```bash
rapela prayer --list
```


# Resources
- https://go.dev/learn/
- https://spf13.com/presentation/building-an-awesome-cli-app-in-go-oscon/
