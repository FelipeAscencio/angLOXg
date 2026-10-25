---
layout: default
title: Diferencias de implementación y agregados
nav_order: 9
---

# 🧬 Diferencias de implementación y agregados

## Extras: agregados al lenguaje

Decidimos ir un paso más allá de los requisitos básicos y extender el lenguaje con dos cosas que Lox estándar no tiene. Las dos están soportadas en los **dos motores de ejecución** y tienen sus propias pruebas.

### Arrays

Nuestro agregado principal. Lox no tiene ningún tipo de dato compuesto, así que un programa no puede guardar una colección de valores sin inventar una clase para cada caso. Agregamos arrays con literales y lectura por posición:

```javascript
var dias = ["lunes", "martes", "miercoles"];
for (var i = 0; i < 3; i = i + 1) print dias[i];

var matriz = [[1, 2, 3], [4, 5, 6]];
print matriz[1][2];   // 6

print [1, "dos", true, nil];
```

Los arrays son **valores de primera clase**: se guardan en variables, se pasan y se devuelven de funciones, viven en propiedades de instancias y se anidan entre sí.

```javascript
class Caja {
  init(cosas) { this.cosas = cosas; }
  primera() { return this.cosas[0]; }
}

print Caja([42, 43]).primera();   // 42
```

**Alcance deliberado.** Esta versión es de **solo lectura**: no hay asignación por índice (`a[0] = 1`), ni largo, ni métodos para agregar o quitar elementos. Fue una decisión consciente de mantener el agregado chico y bien probado en lugar de grande y a medias. Intentar asignar a una posición da un error de sintaxis claro (`destino de asignación inválido`).

**Decisiones de diseño que tomamos:**

*   **El indexado liga al nivel de la llamada**, igual que el `(` de una invocación y el `.` de una propiedad. Por eso `a[0] + 1` suma el elemento y no indexa la suma, y `a[0].x`, `p.lista[0]` y `a[0]()` se combinan en cualquier orden.
*   **Los elementos se evalúan al crear el array**, no al leerlos. En `var a = [x, x * 2]`, cambiar `x` después no modifica lo que guarda `a`.
*   **El índice tiene que ser un entero.** Como todos los números de Lox son de punto flotante, `a[1.5]` no apunta a ninguna posición: preferimos cortar con *"El índice debe ser un número entero"* antes que truncar en silencio.
*   **Un índice fuera de rango es un error de ejecución**, igual que en la mayoría de los lenguajes, en lugar de devolver `nil` y esconder el problema. Cubre también el índice negativo y el array vacío.
*   **La identidad es por referencia**, como en las instancias de clase: `[1] == [1]` da `false`, pero `var a = [1]; var b = a;` hace que `a == b` dé `true`.

**Cómo lo resuelve cada motor:**

| | Intérprete de árbol | Máquina virtual |
|:---|:---|:---|
| **Representación** | Un tipo `LoxArray` con un slice de valores | Un tipo `Array` análogo, en el paquete del compilador |
| **Creación** | El evaluador recorre los elementos del nodo y arma el slice | El opcode `OpArray` lleva como operando la **cantidad de elementos**, que ya están apilados, y se los lleva de una |
| **Lectura** | El nodo `Index` evalúa objeto e índice y valida el rango | El opcode `OpGetIndex` saca índice y array de la pila y valida igual |

El operando de un solo byte de `OpArray` impone un techo de 255 elementos por literal, el mismo límite que ya tenían los argumentos de una función. El compilador avisa si se pasa.

### Operador de módulo (`%`)

Lox estándar no tiene operador de resto. Lo agregamos con la misma precedencia que la multiplicación y la división.

Tuvo una complicación propia de Go que vale contar: el operador `%` de Go **solo acepta enteros**, y los números de Lox son todos de punto flotante. Los dos motores lo resuelven con una conversión explícita a entero, no con el operador nativo.

## Extras: herramientas de desarrollo

*   **Intérprete web (Playground vía WebAssembly):** llevamos angLOXg al navegador. Gracias a la capacidad de Go para compilar a WebAssembly (WASM), el intérprete corre 100% del lado del cliente sin necesitar un servidor backend. Armamos una interfaz web con un editor de scripts completo y una consola interactiva (REPL) en vivo.
*   **Flags de inspección (Modos de ejecución):** Tanto en la consola tradicional (CLI) como en la página web, agregamos la posibilidad de "pausar" la ejecución en distintas etapas del intérprete. Esto es ideal para entender qué pasa por debajo del capó. Implementamos 5 modos visuales:
    1.  **Ejecución completa:** El comportamiento normal de Lox, con el intérprete de recorrido de árbol.
    2.  **Modo escáner (`--escaneo`):** Muestra únicamente la lista de Tokens generada.
    3.  **Modo parser (`--arbol`):** Muestra la estructura del Árbol de Sintaxis Abstracta (AST) formateada como texto.
    4.  **Modo resolver (`--resolucion`):** Ejecuta hasta la fase de resolución de ámbitos y referencias, frenando antes de evaluar el código.
    5.  **Modo máquina virtual (`--vm`):** Compila el programa a bytecode y lo ejecuta en nuestra VM, en lugar de recorrer el árbol.

## El impacto del lenguaje: Go vs. Python

La implementación de referencia de la cátedra (`plox`) está escrita en Python. Al pasar toda esa lógica a **Go**, nos topamos con diferencias arquitectónicas enormes debido a la naturaleza de cada lenguaje. Estas fueron las principales decisiones de diseño que tuvimos que adaptar:

### Tipado estático vs. tipado dinámico

*   *En Python:* Lox es dinámico y Python también, por lo que guardar un número, un string o un booleano en una variable de Python es directo.
*   *En Go:* Al ser de tipado estático estricto, tuvimos que representar los tipos dinámicos de Lox usando el tipo `any` (o `interface{}`). Esto nos obligó a usar *type assertions* (`switch v := valor.(type)`) constantemente en el evaluador para verificar si estábamos sumando dos números o concatenando dos strings, lo que hace al código de Go un poco más verboso pero mucho más seguro en tiempo de compilación.

### Representación del árbol: structs con interfaces marcadoras

*   *En Python:* el AST se construye usando herencia clásica, con una clase base `Expr` vacía y subclases como `BinaryExpr` o `GroupingExpr`. El recorrido no usa el patrón Visitor del libro: aprovecha `functools.singledispatchmethod`, un decorador que define varias versiones del mismo método y deja que Python elija cuál llamar según el tipo del argumento en tiempo de ejecución.
*   *En Go:* como no hay clases ni herencia, modelamos los nodos como `structs` agrupados bajo dos **interfaces marcadoras**, `Expr` y `Stmt`. Cada una declara un único método privado (`isExpr()` / `isStmt()`) que no hace nada: su único propósito es que ningún tipo de afuera del paquete pueda hacerse pasar por un nodo del árbol.

    El recorrido **no usa el patrón Visitor**, sino un `switch` de tipo:

    ```go
    switch n := expresion.(type) {
    case *Binary:
        // ...
    case *Unary:
        // ...
    }
    ```

    En el fondo las dos soluciones hacen lo mismo, despachar según el tipo del nodo, pero Python lo resuelve con su maquinaria de dispatch y en Go hay que escribirlo a mano: no existe nada equivalente a `singledispatchmethod`.

    También descartamos el patrón Visitor del libro, y por una razón concreta: **Go no permite parámetros de tipo en los métodos**. Un visitante genérico necesitaría algo como `Aceptar[T](v Visitante[T]) T`, y eso no compila. La única forma de escribirlo en Go es que todos los métodos devuelvan `any`, obligando a un *type assertion* en cada uso.

    Con el `switch` de tipo, en cambio, cada recorrido devuelve el tipo que le conviene: el que imprime el árbol devuelve `string`, el intérprete devuelve `(any, error)` y el resolvedor devuelve solo `error`. También nos dejó agregar nodos nuevos (los de POO, por ejemplo) sin tocar ninguna interfaz.

    **La contrapartida:** el `switch` de tipo no avisa en tiempo de compilación si a un recorrido le falta un nodo, mientras que el Visitor sí lo haría. Lo compensamos con dos redes: cada `switch` cierra con un `default` que entra en pánico informando el tipo del nodo no contemplado, y tenemos un test que pasa una instancia de cada nodo del paquete por el representador para detectar el olvido antes de que llegue a producción.

### Manejo de errores: tres mecanismos distintos

Go no tiene excepciones, así que cada fase del intérprete resuelve sus errores de la forma que mejor le queda. Terminamos usando tres estrategias diferentes:

*   **En el parser, pánico controlado.** Un descenso recursivo tiene unas quince funciones encadenadas, y devolver errores obligaría a escribir `if err != nil` en todas. Usamos `panic` con un tipo de error privado y `recover` en el borde de cada declaración, que es justamente donde hay que sincronizar para seguir leyendo después de una falla. El pánico nunca escapa del paquete: `Parsear` devuelve `([]Stmt, []error)` como cualquier función de Go.
*   **En el intérprete, errores devueltos.** Acá sí propagamos con `(valor, err)` y chequeamos en cada paso, que es la forma idiomática.
*   **Para el `return` de Lox, un error que no es un error.** Un `return` dentro de funciones anidadas tiene que cortar varios niveles de la pila de golpe. `plox` lo resuelve con una excepción: define `class ReturnValue(Exception)` y la lanza. Nosotros definimos un tipo `ValorRetorno` que **implementa la interfaz `error`** aunque no represente ninguna falla, y lo propagamos por el mismo camino que los errores reales hasta que la llamada a la función lo intercepta. Es el mismo truco, escrito con las herramientas de Go.

## Fase 2: las dos máquinas no se parecen

La segunda entrega agregó un compilador a bytecode y una máquina virtual de pila. Lo interesante es que **no comparten casi nada** con el intérprete de árbol, más allá del escáner.

### El compilador no usa el árbol de sintaxis

El intérprete de recorrido de árbol trabaja en tres pasadas: el parser arma el AST, el resolvedor lo anota con las distancias de los ámbitos, y el evaluador lo recorre.

El compilador de bytecode, en cambio, es de **una sola pasada**: toma la lista de tokens directamente y emite instrucciones mientras lee, sin construir ningún árbol intermedio. Usa un *Pratt parser* con una tabla de reglas por token en lugar de las funciones encadenadas por precedencia del otro parser.

Eso significa que los paquetes `parser`, `sintaxis` y `semantica` **no participan** del modo `--vm`. Las dos implementaciones son independientes de verdad, y es lo que nos permitió compararlas de forma honesta en los benchmarks.

### La misma semántica resuelta de dos formas

Los mecanismos de POO dejaron esto especialmente a la vista. Las dos máquinas hacen lo mismo de cara al usuario, pero por caminos distintos:

| | Intérprete de árbol | Máquina virtual |
|:---|:---|:---|
| **Búsqueda de métodos heredados** | Cada clase guarda un puntero a su superclase y la búsqueda recorre la cadena hacia arriba en cada llamada | Al definir la clase, el opcode `OpInherit` **copia** los métodos de la superclase a la subclase; después se cuelgan los propios, que sobrescriben a los repetidos |
| **Resolución de `this`** | Un entorno propio creado al ligar el método a su instancia | El **slot cero** del marco de llamada, donde la máquina ya deja el receptor: `this` se resuelve como cualquier variable local |
| **Resolución de `super`** | Un entorno con la superclase, capturado por el closure del método | La superclase queda como variable local del ámbito de la clase y cada método la captura como **upvalue** |

El caso que obliga a hacerlo bien es la herencia de varios niveles. Con `class C < B < A`, cuando un método de `B` llama a `super.m()`, tiene que llegar a `A` aunque la instancia sea de `C`. Si `super` se resolviera mirando la clase de la instancia, `C` llevaría de nuevo a `B` y la llamada se volvería recursiva para siempre. Por eso en los dos motores `super` se resuelve **léxicamente**, por el lugar donde está escrito, y no por el tipo del receptor.

Los arrays que agregamos al lenguaje siguieron el mismo patrón: un tipo con un slice en el intérprete de árbol, y un par de opcodes que operan sobre la pila en la máquina virtual.
