package main

import (
    "fmt"
    "net"
    "time"
    "encoding/gob"
)

type Proceso struct {
    Id uint64
    I uint64
    Activo bool
    Write bool
    ActivoC chan bool
}

func (p * Proceso) Start() {
    go func ()  {
        for {
            p.Write = <- p.ActivoC
        }
    }()

    for p.Activo {
        if p.Write {
            fmt.Printf("id %d: %d\n", p.Id, p.I)
        }
        p.I = p.I + 1
        time.Sleep(time.Millisecond * 500)
    }
}

func (p * Proceso) Stop() {
    p.Activo = false
    close(p.ActivoC)
}


func cliente(prcs []*Proceso) {
    var proceso Proceso
    c, err := net.Dial("tcp", ":9999")
    if err != nil {
        fmt.Println(err)
        return
    }
    //err = gob.NewEncoder(c).Encode(proceso)
    err = gob.NewDecoder(c).Decode(&proceso)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Println("PROCESO: ", proceso)
        proceso.Write = true
        go proceso.Start()
        prcs = append(prcs, &proceso)
    }
    c.Close()
}

func main()  {
    prcs := []*Proceso{}
    go cliente(prcs)

    var input string
    fmt.Scanln(&input)
}
