(defun main (a) 
	(print 
		(assert (+ 1 2) 3)
	)
	(print 
		(assert (- 10 2) 8)
	)
	(print 
		(assert (* 12 12) 144)
	)
	(print 
		(assert (/ 100 50) 2)
	)
	(print 
		(assert (str 100) "100")
	)
	(print 
		(assert (int "100") 100)
	)
	(print 
		(assert (float "100.0") 100.0)
	)
	(print 
		(assert (int 100.0) 100)
	)
	(print 
		(assert (str "100.0") "100.0")
	)
	(print 
		(assert (float 100) 100.0)
	)
)
 