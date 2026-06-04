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
rapela bible --list-versions or -l
```

To set a bible version:
```bash
rapela bible --set=<name-of-version> or -s=<name-of-version>
```

# How To Install
1. Download the desired release for your operating system and architecture:<br>
> Go to the [`releases`](https://github.com/kmdinake/rapela/releases/tag/v0.1.0) tab of this repository

2. Extract the folder and place it in your desired location:<br/>
> For instance in `C:\Program Files\` such that you have target folder path of `C:\Program Files\rapela_Windows_x86_64`

3. Append the target folder path to your system `PATH` variable.<br>
For instance on Windows, execute the following command:
```bash
> setx PATH "%PATH%;C:\Program Files\rapela_Windows_x86_64"
```
Thereafter, refresh your terminal to get the latest environment variables.

4. Verify installation
Execute the following command to see the verse of the day
```bash
> rapela
```
Then you will see the verse of the day:<br>
```bash
Dumela ngwana waka! Tseya sebaka se go rapela.
Verse Of The Day: (tn-olef) Johane II 1:12 - Ke ne ke eletsa go bua tse dingwe, mme ga ke batle go di bua mo lokwalong lo, gonne ke solofela go tla go lo bona ka bofefo mme ke gone re tlaa buisanyang ka tsone re bo re nna le nako ya boitumelo.
```

## Troubleshoot Installation
1. `Access is denied` on execution of `rapela` Windows
> Assign your user profile full-control to the target folder. <br>If this doesn't resolve the issue, and you have an anti-virus software installed.<br> Then add an exclusion for the target folder. 

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
