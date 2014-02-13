with Ada.Text_IO;
with Ada.Strings.Unbounded;

procedure Test_Task is
	package SUB renames Ada.Strings.Unbounded;

	type Task_Data is
		record
			DelayTime : Standard.Duration;
			Message : SUB.Unbounded_String;
		end record;
	
	task type Print_Task (D : access Task_Data);

	task body Print_Task is
	begin
		loop
			Ada.Text_IO.Put_Line(SUB.To_String(D.Message));
			delay D.DelayTime;
		end loop;
	end;

	Data_TaskA : aliased Task_Data := (DelayTime => 1.0, Message => SUB.To_Unbounded_String("Hello"));
	Data_TaskB : aliased Task_Data := (DelayTime => 2.0, Message => SUB.To_Unbounded_String("World"));
	A : Print_Task(Data_TaskA'Access);
	B : Print_Task(Data_TaskB'Access);
begin
	null;
end Test_Task;
