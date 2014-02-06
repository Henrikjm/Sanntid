/*We want to:
Read buttons
Read floor indicators
Set lights
Set motor speed
*/

package main

//#include <stdlib.h>
//#include "io.h"
//#include "channels.h"
import "C"

import (
	"fmt"
)

func elev_init() {
	if !C.io_init() {
		return 0
	}


	for i := 0; i < C.N_FLOORS; i++ {
    if (i != 0){
    	
    }

    if i != N_FLOORS-1 {
    	
    }

	}

    C.io_clear_bit(C.LIGHT_STOP)
    C.io_clear_bit(C.DOOR_OPEN)
    return 1
}


int elev_init(void){
    // Init hardware
    if (!io_init())
        return 0;

    // Zero all floor button lamps
    for (int i = 0; i < N_FLOORS; ++i) {
        if (i != 0)
            elev_set_button_lamp(BUTTON_CALL_DOWN, i, 0);

        if (i != N_FLOORS-1)
            elev_set_button_lamp(BUTTON_CALL_UP, i, 0);

        elev_set_button_lamp(BUTTON_COMMAND, i, 0);
    }

    // Clear stop lamp, door open lamp, and set floor indicator to ground floor.
    elev_set_stop_lamp(0);
    elev_set_door_open_lamp(0);
    elev_set_floor_indicator(0);

    // Return success.
    return 1;
}

func main() {
	fmt.Println("test")
	fmt.Println(C.abs(-5))
	fmt.Println("jøss")
	fmt.Println("Se her da:", C.PORT3)
}
