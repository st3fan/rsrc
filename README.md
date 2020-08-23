# RSRC - Macintosh Resource File Parser

*Stefan Arentz, March 2018*

This package let's you parse classic MacOS resource files. Because, reasons ...

```
import (
    "fmt"

    "github.com/st3fan/rsrc"
)

func main() {
    file, err := rsrc.FromPath("SolarianII")
    if err != nil {
        panic(err)
    }

    for i := 0; i < file.CountResources("snd "); i++ {
        if r, ok := file.GetResource("snd ", i); ok {
            fmt.Printf("Found a sound resource; ID=%d Size=%d Name=%s\n", r.ID, len(r.Data), r.Name)
        }
    }
}
```

```
$ go run example.go
Found a sound resource; ID=5016 Size=4493 Name=Present Bounce
Found a sound resource; ID=5011 Size=4124 Name=Ship Death
Found a sound resource; ID=5001 Size=456 Name=Nasty Death
Found a sound resource; ID=5014 Size=1707 Name=Congas (Buddhabuddha)
Found a sound resource; ID=5015 Size=5145 Name=Synth Twang
Found a sound resource; ID=5007 Size=5660 Name=Yes.
Found a sound resource; ID=5004 Size=3269 Name=Gareth Yeah
Found a sound resource; ID=5009 Size=10531 Name=Hey, Hey, Hey!
Found a sound resource; ID=5006 Size=15211 Name=Missile Explosion
Found a sound resource; ID=5005 Size=15852 Name=Missile Launch
Found a sound resource; ID=5008 Size=12199 Name=Laugh 2
Found a sound resource; ID=5013 Size=7878 Name=Essence du Tim
Found a sound resource; ID=5003 Size=4853 Name=Ben Oing
Found a sound resource; ID=5017 Size=21071 Name=Explosion
Found a sound resource; ID=5018 Size=922 Name=Roosters In Gaffa
Found a sound resource; ID=5002 Size=7237 Name=Supply Ship Alert
Found a sound resource; ID=5010 Size=24716 Name=Fun With Trumpets
Found a sound resource; ID=5000 Size=12384 Name=Cash Register
Found a sound resource; ID=5019 Size=18380 Name=Laugh 1
Found a sound resource; ID=5020 Size=20844 Name=Scissor Snip
Found a sound resource; ID=5022 Size=12092 Name=Grunt Of Life
Found a sound resource; ID=5021 Size=34924 Name=Mary Dismembered & Screaming
Found a sound resource; ID=5023 Size=8736 Name=Damn!
Found a sound resource; ID=5025 Size=3148 Name=Keir's A Drip
Found a sound resource; ID=5026 Size=4828 Name=Keir's A Double Drip
Found a sound resource; ID=5024 Size=37073 Name=Shield Hit
Found a sound resource; ID=5027 Size=67805 Name=More Fun With Trumpets
Found a sound resource; ID=5012 Size=46276 Name=Andrei Praise Jesus
Found a sound resource; ID=5028 Size=6964 Name=Andrei You Lose
```
