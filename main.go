package main
import (
    "github.com/roymiloh/AsciiArt"
    "io"
    "os"
    "bufio"
)

func main() {
    inputFile, err := os.Open("src/github.com/roymiloh/AsciiArtExample/data/input.txt")
    if err != nil {
        panic(err)
    }
    outputFile, err := os.Create("src/github.com/roymiloh/AsciiArtExample/data/output.txt")
    if err != nil {
        panic(err)
    }

    r := AsciiArt.NewReader(inputFile)
    wr := bufio.NewWriter(outputFile)

    // a channel indicates on finish of writing process
    done := make(chan bool, 1)
    flushRate := 10

    // buffered channel for throttling purposes
    chRec := make(chan *AsciiArt.Record, flushRate)

    // writing process in his own goroutine
    go func() {
        recCount := 0

        for record := range chRec {
            recCount += 1
            record.Write(wr)

            if recCount % flushRate == 0 {
                wr.Flush()
            }
        }

        wr.Flush()
        done <- true
    }()

    // reading and sending values to chRec
    for {
        record, err := r.Read()
        if err != nil {
            if err == io.EOF {
                close(chRec)
                break;
            } else {
                // terminate the process if anything unexpected happens
                // ..or do something else, it actually depends on the error.
                panic(err)
            }
        }
        chRec <- record
    }

    // waiting...
    <- done
}
