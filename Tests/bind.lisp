(defun main (a)
        (bind (v1 v2)
        (test a))
	(print v1 v2)
)
(defun test (a)
        (setq i1 10)
        (setq i2 1)
	(values i1 i2)
)
 