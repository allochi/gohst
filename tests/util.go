package main

import (
	"allochi/gohst"
	. "allochi/gohst/plays/models"
	"time"
)

func cleanup() {
	Contactizer, _ := gohst.GetDataStore("Contactizer")
	Contactizer.Drop(Greek{}, true)
}

func greekAlphabet() []Greek {

	return []Greek{
		Greek{0, "Αα Alpha", time.Now(), time.Now()},
		Greek{0, "Ββ Beta", time.Now(), time.Now()},
		Greek{0, "Γγ Gamma", time.Now(), time.Now()},
		Greek{0, "Δδ Delta", time.Now(), time.Now()},
		Greek{0, "Εε Epsilon", time.Now(), time.Now()},
		Greek{0, "Ζζ Zeta", time.Now(), time.Now()},
		Greek{0, "Ηη Eta", time.Now(), time.Now()},
		Greek{0, "Θθ Theta", time.Now(), time.Now()},
		Greek{0, "Ιι Iota", time.Now(), time.Now()},
		Greek{0, "Κκ Kappa", time.Now(), time.Now()},
		Greek{0, "Λλ Lambda", time.Now(), time.Now()},
		Greek{0, "Μμ Mu", time.Now(), time.Now()},
		Greek{0, "Νν Nu", time.Now(), time.Now()},
		Greek{0, "Ξξ Xi", time.Now(), time.Now()},
		Greek{0, "Οο Omicron", time.Now(), time.Now()},
		Greek{0, "Ππ Pi", time.Now(), time.Now()},
		Greek{0, "Ρρ Rho", time.Now(), time.Now()},
		Greek{0, "Σσ Sigma", time.Now(), time.Now()},
		Greek{0, "Ττ Tau", time.Now(), time.Now()},
		Greek{0, "Υυ Upsilon", time.Now(), time.Now()},
		Greek{0, "Φφ Phi", time.Now(), time.Now()},
		Greek{0, "Χχ Chi", time.Now(), time.Now()},
		Greek{0, "Ψψ Psi", time.Now(), time.Now()},
		Greek{0, "Ωω Omega", time.Now(), time.Now()},
	}

}
