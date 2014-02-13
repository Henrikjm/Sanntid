/* 
p2pNetwork implements a robust peer-to-peer network over TCP in a local area network. 

A single network node establishes its connection by broadcasting on UDP its own listen address, then waiting for connections from the other nodes on the network. After this a consistency check is performed on the initating node and on all nodes on the existing network. If all nodes agree on the new network connection to the existing peer-to-peer network, the new node is added.

This sequence diagram illustrates how a new network node connects to an existing peer-to-peer network: https://www.lucidchart.com/documents/view/4bfb-aa94-516aac93-b8ee-271f0a0040c6
*/
package p2pNetwork
