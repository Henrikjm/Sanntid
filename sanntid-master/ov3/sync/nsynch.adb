with Ada.Text_IO;

procedure NSynch is
	N : constant := 10;
	TaskDelay : constant := 2.0;

	protected Manager is
		entry Synchronize;
	private
		Final : Boolean := False;
		entry Wait;
	end Manager;

	protected body Manager is
		entry Synchronize when True is
		begin
			Final := Wait'Count = N-1;
			requeue Wait;
		end Synchronize;

		entry Wait when Final is
		begin
			null;
		end Wait;
	end Manager;

	task type Worker;
	task body Worker is
	begin
		loop
			Ada.Text_IO.Put("!");
			delay TaskDelay;
			Manager.Synchronize;
			Ada.Text_IO.Put(".");
			delay TaskDelay;
			Manager.Synchronize;
		end loop;
	end Worker;

	Workers : array (1 .. N) of Worker;
begin
	null;
end NSynch;






