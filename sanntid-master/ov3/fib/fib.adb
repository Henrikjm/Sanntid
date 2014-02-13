with Ada.Integer_Text_IO;
with Ada.Text_IO;
with Ada.Command_Line;
with Ada.Strings.Fixed;
with Ada.Strings;

procedure Fib is
	function Fibonacci (n : Natural) return Natural is
		A, B, C : Natural;
	begin
		if n = 0 then
			return 0;
		elsif n <= 2 then
			return 1;
		else
			A := 1;
			B := 1;
			FibLoop:
			for I in Integer range 3 .. n loop
					C := A + B;
					A := B;
					B := C;
			end loop FibLoop;
			return C;
		end if;
	end Fibonacci;
	
	package SF renames Ada.Strings.Fixed;
	
	N,F : Natural;
begin
	N := Integer'Value(Ada.Command_Line.Argument(1));
	F := Fibonacci(N);
	Ada.Text_IO.Put_Line(Integer'Image(F));
end Fib;
