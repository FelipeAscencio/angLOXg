package compilador

import "testing"

func TestNuevaFuncionArrancaVacia(t *testing.T) {
	fn := NuevaFuncion()

	if fn.Aridad != 0 {
		t.Errorf("la aridad es %d y se esperaba 0", fn.Aridad)
	}

	if fn.UpvalueCount != 0 {
		t.Errorf("la cuenta de upvalues es %d y se esperaba 0", fn.UpvalueCount)
	}

	if fn.Nombre != "" {
		t.Errorf("el nombre es %q y se esperaba vacío", fn.Nombre)
	}

	if fn.Chunk == nil {
		t.Error("la función nueva no tiene chunk")
	}
}

func TestRepresentacionDeFuncion(t *testing.T) {
	casos := []struct {
		nombre   string
		funcion  *Funcion
		esperado string
	}{
		{
			nombre:   "función anónima",
			funcion:  NuevaFuncion(),
			esperado: "<script>",
		},
		{
			nombre:   "función con nombre",
			funcion:  &Funcion{Nombre: "sumar"},
			esperado: "<fn sumar>",
		},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			if obtenido := caso.funcion.String(); obtenido != caso.esperado {
				t.Errorf("se obtuvo %q y se esperaba %q", obtenido, caso.esperado)
			}
		})
	}
}

func TestNuevoClosureReservaLosUpvalues(t *testing.T) {
	fn := &Funcion{Nombre: "contador", UpvalueCount: 3, Chunk: NuevoChunk()}

	cl := NuevoClosure(fn)

	if cl.Funcion != fn {
		t.Error("el closure no apunta a la función que se le pasó")
	}

	if len(cl.Upvalues) != 3 {
		t.Errorf("reservó %d upvalues y se esperaban 3", len(cl.Upvalues))
	}
}

func TestClosureSinUpvalues(t *testing.T) {
	cl := NuevoClosure(NuevaFuncion())

	if len(cl.Upvalues) != 0 {
		t.Errorf("reservó %d upvalues y se esperaban 0", len(cl.Upvalues))
	}
}
