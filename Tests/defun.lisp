(defun main (a) 
	(print "Begin") 
	(print 
		(test a 7)
	)	 
	(print 
		(test_minus 12 7)
	)	 
	(print "End")
)
(defun (test a b) 
	(print a) 
	(print b) 
	(+ a b)
)

(defun (test_minus a b) 
	(print a) 
	(print b) 
	(- a b)
)
 