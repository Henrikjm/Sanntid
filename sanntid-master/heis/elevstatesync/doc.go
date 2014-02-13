/* 
This modules main task is to syncronize its state with other elevators connected to the local area network. In addition it contains logic to delegate and add or remove incoming orders. 

All elevators recieve and register all state changes via syncronization messages. Since we expect at least one elevator to always run, this ensures that no orders are lost. This sequence diagram illustrates how a single state change is done
https://www.lucidchart.com/documents/view/4ff9-c3a8-516aad58-825d-25bc0a009224

A call to SyncXXX() blocks until the message is sent and returns a copy of the updated state to the caller. If a message is recieved, the state is updated and syncronized on the NewStateChan of the SyncState object.

This module uses the p2pNetwork module to establish the TCP connections in a peer-to-peer fashion.
*/
package elevstatesync
