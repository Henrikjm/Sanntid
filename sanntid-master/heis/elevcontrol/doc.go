/* 
elevcontrol contains the main functionality to control a single elevator
The module generates events to a finite state machine by listening to events from
the elevdriver module. State changes are syncronized with the elevstatesync module

Take a look at the finite state machine diagram at 
https://www.lucidchart.com/documents/view/4cb2-73e8-516aa800-96f0-56930a0041c5
*/
package elevcontrol
