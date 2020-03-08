package main

type ONB struct {
	u, v, w Tuple
}

func (o ONB) local(a Tuple) Tuple {
	return o.u.MulScalar(a.x).Add(o.v.MulScalar(a.y)).Add(o.w.MulScalar(a.z))
}

func buildFromW(n Tuple) ONB {
	var o ONB
	var a Tuple

	o.w = n.Normalize()
	if o.w.x > 0.9 {
		a = Tuple{0, 1, 0, 0}
	} else {
		a = Tuple{1, 0, 0, 0}
	}

	o.v = o.w.Cross(a).Normalize()
	o.u = o.w.Cross(o.v)

	return o
}
