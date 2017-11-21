package semLib

import "fmt"



//==============================================================
// Semáforo
//==============================================================
type Semaphore struct{
    n int
    inc chan int
    dec chan int
}
//Construtor
func SemInit(n int, inc, dec chan int) Semaphore {
	return Semaphore{n, inc, dec}
}
func MutexInit(inc, dec chan int) Semaphore {
	return Semaphore{1, inc, dec}
}
func (s Semaphore) Signal(n int) {
    s.inc<-n
}
func (s Semaphore) Wait(n int) {
    s.dec<-n
}
func (s Semaphore) Run() {
    for {
        if s.n == 0 {
            <-s.inc
            s.n++
        } else if s.n > 0 {
            select{
            case <-s.inc:
                s.n++
            case <-s.dec:
                s.n--
            }
        } else {
            fmt.Println("Erro: semaforo menor que zero")
            return
        }
    }
}
//=====================================================================