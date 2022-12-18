# Reef Controller

## Aquarium automatic water top off controller

Components
* [Non-contact Electronic Water Level Sensor](https://www.amazon.com/gp/product/B07Z64CSLQ/ref=ppx_yo_dt_b_search_asin_title?ie=UTF8&psc=1)
  * [datasheet](docs/Taidacent%20Mini%20External%20Sticker%20Intelligent%20Non-contact%20Electronic%20Water%20Level%20Sensor%20High%20Low%20Level%20Output%20Liquid%20Level%20Sensor%20Switch%20for%20Water%20Tank%20Fish%20Tank%20.pdf)
* [Float Switch](https://www.amazon.com/gp/product/B072QCHQ2P/ref=ppx_yo_dt_b_search_asin_title?ie=UTF8&psc=1)
* [Bayite 12V DC Water Pressure Diaphragm Pump](https://www.amazon.com/gp/product/B01N75ZIXF/ref=ppx_yo_dt_b_search_asin_title?ie=UTF8&psc=1)
  * Vertical Suction Lift: 5ft (1.5m). Adjustable CUT-OFF Pressure : default 80 PSI, MAX 100 PSI. Amp Draw: 3.0AMP.It doesn't mean 4L/min jet at 80 PSI.
* Water Reservoir
  * Any water container/bucket that you wish to use. I'm using a 5 gallon bucket with a water tight seal to prevent evaporation, with Aragonite sand to keep a consistent 8.2 ph.
* [MT3608 DC-DC Step Up Power Booster](https://www.amazon.com/gp/product/B089JYBF25/ref=ppx_yo_dt_b_search_asin_title?ie=UTF8&psc=1)
  * [datasheet](docs/MT3608.pdf)
* [3v Relay](https://www.amazon.com/gp/product/B08W3XDNGK/ref=ppx_yo_dt_b_search_asin_title?ie=UTF8&psc=1)
  * [5v Relays](https://www.amazon.com/gp/product/B095YD3732/ref=ppx_yo_dt_b_search_asin_title?ie=UTF8&psc=1) if you're using a 5v system, e. g. Arduino uno
* Resistors
  * 3x 330Ω resistors
  * 1x 10kΩ resistor
* LEDs
  * 1x green
  * 1x yellow
  * 1x red
* JST PH2.0 connectors
  * 1x 3 pin for water level sensor
  * 1x 2 pin for water pump

### Circuit Diagrams


[Circuit Diagram Source](https://crcit.net/c/d85ce8abab73407c9d194f981ed2c907)<br>
![Basic Reef Controller Circuit Diagram](docs/images/Reef%20Controller%20Circuit%202.png)



### PCB Prototype Pinout

* Water level sensor
  * Power: D
  * Data: C
  * Ground: E
