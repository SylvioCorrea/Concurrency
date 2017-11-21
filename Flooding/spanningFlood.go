//Sylvio Correa e Annderson Oreto

package main

import (
        "fmt"
        "math/rand"
        "time"
        )
//

// Struct que define uma mensagem
type Message struct {
    senderID int
    messageID int
}

// Struct que guarda um canal juntamente com a identidade do nodo
// ao qual esse canal leva
type EdgeChannel struct {
    ownerID int
    edgeC chan Message
}

// Struct que define a estrutura de um nodo do grafo
type NodeProc struct {
    id int
    receive chan Message
    edges []EdgeChannel
}

// O nodo executa 3 processos ao mesmo tempo:
// O primeiro é o mutex que regula o acesso dos outros dois processos
// ao arquivo de mensagens do nodo. O segundo é o processo que cuida
// do recebimento de mensagens, o terceiro é o que cuida do envio de mensagens.
func (np *NodeProc) runNodeProc() {
    go np.receiverProtocol()
    go np.senderProtocol()
}

// Processo do nodo que cuida do recebimnto de mensagens
func (np *NodeProc) receiverProtocol() {
    for {
         msg := <- np.receive
         fmt.Printf("O nodo %d recebe a mensagem %d do nodo %d e inicia um processo de flood\n", np.id, msg.messageID, msg.senderID)
         
         //Altera a identidade do mensageiro no struct mensagem
         //a ser propagado no grafo.
         previousSender := msg.senderID
         msg.senderID = np.id
         
         //Inicia o processo de flood em paralelo para que o processo de
         //recebimento possa continuar, evitando o deadlock.
         go np.flood(msg, previousSender)
    }
}

// Processo do nodo que cuida do envio de NOVAS mensagens para o grafo 
func (np *NodeProc) senderProtocol() {
    for {
        //O nodo aguarda algum tempo entre 1 e 3 segundos e então
        //envia uma mensagem aleatória.
        waitTime := (rand.Int() % 3) + 1
        time.Sleep(time.Duration(waitTime) * time.Second)
        
        msg := Message{np.id, rand.Int() % 1000}
        fmt.Printf("O nodo %d inicia o envio da mensagem %d\n", np.id, msg.messageID)
        np.floodAll(msg)
    }
}

//=============================
// Processos de flood:
//=============================
// Flood que envia mensagens para todos os nodos adjacentes ao remetente,
// menos o nodo do qual ele recebeu a mensagem.
func (np *NodeProc) flood(msg Message, previousSender int) {
    for i := range np.edges {
        if np.edges[i].ownerID != previousSender {
            fmt.Println(np.id, "mandando para", np.edges[i].ownerID)
            np.edges[i].edgeC <- msg
        }
    }
    
}
// Flood que envia mensagens para todos os nodos adjacentes ao remetente,
// sem exceção. Usado quando o remetente é o originador da mensagem.
func (np *NodeProc) floodAll(msg Message) {
    for i := range np.edges {
        fmt.Println(np.id, "mandando para", np.edges[i].ownerID)
        np.edges[i].edgeC <- msg
    }
    
}






//===========================================================





func main() {
    nProcs := 5
    
    //Matriz que define ligações para cada um dos nodos.
    //É fundamental que os dados da matriz descrevam uma
    //árvore de espalhamento.
    connections := [][]int{
                           []int{1,4},    //Nodo 0 está ligado a Nodos 1 e 4
                           []int{0,2,3},  //Nodo 1...
                           []int{1},      //Nodo 2...
                           []int{1},      //Nodo 3...
                           []int{0},    //Nodo 4...
                           }
    
    //Cria um slice de structs que guardam um canal de mensagens e a
    //identidade do dono do canal.
    nodeChannels := make([]EdgeChannel, nProcs)
    for i:= 0; i<nProcs; i++ {
        nodeChannels[i] = EdgeChannel{i, make(chan Message)}
    }
    
    //Laço que constrói e executa cada um dos nodos.
    for i:= 0; i<nProcs; i++ {
        
        //Agrupa os canais de cada um dos nodos ligados ao nodo i.
        edgesTemp := make([]EdgeChannel, len(connections[i]))
        for j := range connections[i] {
            edgesTemp[j] = nodeChannels[connections[i][j]]
        }
        
        np := NodeProc{
                       id: i,
                       receive: nodeChannels[i].edgeC,
                       edges: edgesTemp,
                       }
        
        np.runNodeProc()
    }
    
    //Apenas um jeito tosco de interromper a main:
    quit := make(chan int)
    <-quit
}


























