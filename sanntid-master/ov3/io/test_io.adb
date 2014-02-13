with Ada.Text_IO;
use Ada.Text_IO;
with Ada.Strings;
use Ada.Strings;

procedure Test_IO is
	MaxLength : Natural := 16;
	Str : String (1..MaxLength);
	Last : Natural;
begin
	SuperLoop:
	loop
		Ada.Text_IO.Get_Line(Str, Last);
		if Last > MaxLength then
			Ada.Text_IO.Put_Line("Overflow!");
			raise Ada.Strings.Length_Error;
		elsif Last = 0 then
			Ada.Text_IO.Put_Line("Empty string not allowed");
			raise Ada.Strings.Length_Error;
		end if;
		Ada.Text_IO.Put_Line(Str(1..Last));
	end loop SuperLoop;
end Test_IO;
