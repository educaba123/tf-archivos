# tf-archivos
tf archivos
## Descripción del problema y motivación

La problemática para el caso abordado se basa en la dificultad de los concesionarios de vehículos para la venta de sus productos y de los asesores comerciales para concretar ventas ante la negativa de varios clientes. Por la falta de eficiencia y con el fin de maximizar las ventas y ganancias, se pueden establecer comisiones justas que motiven a los asesores a promocionar los autos para que los clientes decanten con su adquisición y, con ello, el vendedor encargado reciba una comisión del vehículo que vendió. Ante este problema se busca optimizar las estrategias de venta, como ajustar las tasas de comisión para incentivar a los asesores a vender vehículos costosos o liquidar autos poco demandados. Además, se podrían planificar mejor los presupuestos y estrategias financieras al realizar predicciones de mayor precisión respecto a los ingresos y a las comisiones a pagar.

La motivación para el desarrollo de este trabajo se basa en la mejora de toma de decisiones de los concesionarios, ya que con el uso de datos precisos y análisis predictivo, se puede ajustar el precio de los vehículos y las comisiones para maximizar las ganancias. Asimismo, el ajuste de estas estrategias de ventas, basados en el análisis de datos, puede aumentar la competitividad tanto dentro de la empresa como en el mercado. A través de ello, se podría adelantar las tendencias que existen en el mercado de venta de vehículos frente a otros concesionarios. Y por último, un sistema de comisiones optimizado puede incentivar a los asesores comerciales. Eso sucede al estar relacionado el desempeño que realizan con el salario que perciben, pues puede motivar a que trabajen de forma más eficiente.

## Objetivos

- Implementar un sistema distribuido para realizar cálculos de regresión lineal en paralelo utilizando múltiples nodos.
- Aplicar un algoritmo de regresión lineal para encontrar la relación entre dos variables.
- Establecer comunicación entre el coordinador y los nodos trabajadores para distribuir tareas y recolectar resultados.
- Permitir la entrada dinámica de direcciones IP y puertos para nodos y servidor desde la consola.
- Asegurar que el sistema sea escalable y eficiente para manejar grandes volúmenes de datos.
- Manejar la concurrencia y sincronización en un entorno distribuido.

## Sustento del algoritmo utilizado

La regresión lineal es una técnica estadística fundamental que se utiliza para modelar la relación entre una variable dependiente y una o más variables independientes. Según Dagnino (2014), la regresión lineal se utiliza para identificar relaciones potencialmente causales entre variables o, en contextos donde la relación causal es clara, para predecir el valor de una variable en función de otra. La fórmula de la regresión lineal es:


donde:
- y es la variable dependiente (en este caso, “sale_price”).
- x es la variable independiente (en este caso, “commission_rate”).
- m es la pendiente de la línea de regresión.
- b es la intersección con el eje y.

Por medio del algoritmo de Regresión Lineal podemos realizar un cálculo distribuido que nos permite manejar grandes volúmenes de datos de forma eficiente, además que se adapta a la implementación de más nodos para una mayor escalabilidad del sistema. Por ello, este algoritmo tiene la capacidad de realizar predicciones precisas.

## Explicación del dataset y estructura de datos

El dataset utilizado contiene información de un registro de ventas de vehículos, cada uno de ellos representa una transacción de venta de un automóvil. La estructura de datos es la siguiente:

- **Date**: Fecha de la venta
- **Salesperson**: Nombre del vendedor que realizó la venta
- **Customer Name**: Nombre del cliente que adquirió el vehículo
- **Customer Dni**: Número de DNI del cliente
- **Car Make**: Marca del automóvil vendido
- **Car Model**: Modelo del automóvil vendido
- **Car Year**: Año de fabricación del automóvil
- **Sale Price**: Precio de venta del automóvil
- **Commission Rate**: Tasa de comisión para la venta
- **Commission Earned**: Comisión ganada por el vendedor

El programa aplicará el algoritmo de regresión lineal para modelar la relación entre el precio de venta (“sale_price”) y la tasa de comisión (“commission_rate”), de forma que se calcularía la pendiente y la intersección de la línea de regresión. Al final el resultado de este algoritmo es enviado al cluster.
