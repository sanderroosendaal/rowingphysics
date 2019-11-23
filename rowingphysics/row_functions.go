package main

import (
  "fmt"
  "math"
)

// Global constants

// water density (kg/m^3)

const rho = 999.97

// maximum lift coefficient for blade
const CLmax = 1.0

func drag_eq(displacement float64, velo float64,
  alfaref float64, doprint int, constantdrag int) float64 {
  if alfaref == 0 {
    alfaref = 3.5
  }
  var corr = alfaref/3.5

  // usual coefficient - boat is a spheroid
  const c1 float64 = 1./3.
  const c2 float64 = 2./3.
  const d1 float64 = 0.06
  const d2 float64 = 28.0
  const d3 float64 = 0.10891
  var a1 float64 = 0
  const a3 float64 = 0

  var beam float64 = a1+d1*math.Pow(displacement,c1)
  var boatlength float64 = d2*beam
  var wetted_area float64 = a3+d3*math.Pow(displacement,c2)

  const kinvis float64 = 1.19e-6

  // var D float64 = displacement/rho
  var Re float64 = boatlength*velo/kinvis
  var Cf float64 = 0.075/(math.Pow((math.Log10(Re)-2.0),2))
  var alpha float64 = 0.5*rho*wetted_area*Cf
  alpha = alpha*corr
  a1 = alpha/0.8

  if constantdrag == 1 {
    a1 = alfaref
  }

  if doprint == 1 {
    fmt.Print("------ Drag resistance data ----------------\n")
    fmt.Printf("Corr : %f\n", corr)
    fmt.Printf("Beam : %f\n", beam)
    fmt.Printf("Boat length : %f\n", boatlength)
    fmt.Printf("Wetted Area : %f\n", wetted_area)
    fmt.Printf("alpha skin : %f\n",alpha)
    fmt.Printf("alpha total : %f\n",a1)
    fmt.Print("------  Drag resistance data ---------------\n")
    fmt.Print("\n")
  }

  var W2 = a1*math.Pow(velo,2)

  return(W2)
}

func vboat(mc float64, mb float64, vc float64) float64 {
  return(mc*vc/(mc+mb))
}

func vhandle(v float64, lin float64, lout float64, mc float64, mb float64) float64 {
  var gamma = mc/(mc+mb)
  var vc = lin*v/(lout+gamma*lin)
  return(vc)
}

func d_recovery(dt, v, vc, dvc, mc, mb, alef, F float64) float64 {
  var dv = 0.0
  var vb = vboat(mc,mb,vc)
  // var dvb = vboat(mc,mb,dvc)
  var Ftot = F-alef*math.Pow((v-vb),2)
  dv = dt*Ftot/(mb+mc)
  return(dv)
}

func d_stroke(dt, v, vc, dvc, mc, mb, alef, F float64) float64 {
  var dv = 0.0
  var vb = vboat(mc,mb,vc)
  var Ftot = F-alef*math.Pow((v-vb),2)

  dv = dt*Ftot/(mb+mc)

  return(dv)
}

func de_footboard(mc,mb,vs1,vs2 float64) float64 {
  var de = 0.0
  var vt = 0.0

  var vb1 = vboat(mc,mb,vs1)
  var vb2 = vboat(mc,mb,vs2)

  var vmb1 = vt-vb1
  var vmb2 = vt-vb2

  var vmc1 = vt+vs1-vb1
  var vmc2 = vt+vs2-vb2

  var e_1 = 0.5*(mb*math.Pow(vmb1,2))+0.5*mc*math.Pow(vmc1,2)
  var e_2 = 0.5*(mb*math.Pow(vmb2,2))+0.5*mc*math.Pow(vmc2,2)
  var e_t = 0.5*(mc+mb)*math.Pow(vt,2)


  // var de2 = e_2 - e_t
  // var de1 = e_1 - e_t
  var samesign bool = (vs1 * vs2) > 0

  if samesign {
    de = e_2 - e_1
    if e_2 - e_1 < 0 {
      de = 0
    }
    } else {
      de = e_2 - e_t
      if de < 0 {
        de = 0
      }
    }

    return(de)
  }

  type rig struct {
    lin float64
    mb float64
    bladelength float64
    lscull float64
    Nrowers int32
    roworscull string
    span float64
    catchangle float64
    dragform float64
    __spread float64
    __bladearea float64
  }

  func (rg *rig) spread() float64 {
    if rg.roworscull == "scull" {
      return (rg.span/2.)
    }
    return rg.__spread
  }

  func (rg *rig) overlap() float64 {
    if rg.roworscull == "scull" {
      return 2.*rg.lin-rg.span
    }
    return rg.lin-rg.__spread
  }

  func (rg *rig) buitenhand() float64 {
    if rg.roworscull == "scull" {
      return rg.span-2.*rg.lin*math.Cos(rg.catchangle)
    }
    return rg.__spread-rg.lin*math.Cos(rg.catchangle)
  }

  func (rg *rig) bladearea() float64 {
    if rg.roworscull == "scull" {
      return 2.*rg.__bladearea
    }
    return rg.__bladearea
  }

  func (rg *rig) dcatch() float64 {
    return(rg.lin*math.Sin(rg.catchangle))
  }

  func (rg *rig) oarangle(x float64) float64 {
    var dist = rg.dcatch()+x
    var angle = math.Asin(dist/rg.lin)
    return(angle)
  }

  func newRig(lin float64,mb float64,lscull float64,
    span float64,spread float64, roworscull string,
    catchangle float64, bladearea float64, bladelength float64,
    Nrowers int32, dragform float64) *rig {
      return &rig{
        lin: lin,
        mb: mb,
        bladelength: bladelength,
        lscull: lscull-0.5*bladelength,
        Nrowers: Nrowers,
        roworscull: roworscull,
        span: span,
        catchangle: catchangle,
        dragform: dragform,
      }
  }

  func blade_force(oarangle float64, rigging rig, vb, fblade float64) {

  }

  func main() {
    var res  = drag_eq(100,4.5,3.5,1,0)
    fmt.Printf("Result : %f\n",res)

    var dv = d_recovery(0.01, 4.5, 1.0, 0.1, 80., 14., 3.5, 100.)
    fmt.Printf("Dv = %f\n",dv)
  }
