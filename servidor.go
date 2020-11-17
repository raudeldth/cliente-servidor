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
    //p.Activo = false
    close(p.ActivoC)
}

func eliminaProceso(prcs []*Proceso, s int) {
    prcs = append(prcs[:s], prcs[s+1:]...)
}

func servidor(prcs []*Proceso)  {
    s, err := net.Listen("tcp", ":9999")
    if err != nil {
        fmt.Println(err)
        return
    }

    for {
        c, err := s.Accept()
        if err != nil {
            fmt.Println(err)
            continue
        }
        go handleClient(c, prcs)
    }
}

func handleClient(c net.Conn, prcs []*Proceso)  {
    prcs[0].Stop()
    err := gob.NewEncoder(c).Encode(prcs[0])
    eliminaProceso(prcs, 0)
    //err := gob.NewDecoder(c).Decode(&proceso)
    if err != nil {
        fmt.Println(err)
        return
    } /*else {
        fmt.Println("Mensaje2: ", prcs[0])
    }*/
}

func main()  {
    prcs := []*Proceso{}
    for i := 0;  i < 5; i++ {
        c := make(chan bool)
        pro := Proceso{uint64(i), 0, true, true, c}
        time.Sleep(time.Millisecond * 1)
        go pro.Start()
        prcs = append(prcs, &pro)
    }

    go servidor(prcs)

    var input string
    fmt.Scanln(&input)
}
