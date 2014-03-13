package main

import(
   "networkIO"
   ."net"
   )


func main(){

  chanCon := make(chan *UDPConn,1)
  go networkIO.ListenToNetwork(chanCon, "inf", "20661")
  networkIO.SendToNetwork(chanCon, "20661", "writeToConsole")
}
