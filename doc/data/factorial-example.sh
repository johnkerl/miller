mlr --opprint --from data/small put '
    func f(n) {
        if (is_numeric(n)) {
            if (n > 0) {
                return n * f(n-1);
            } else {
                return 1;
            }
        }
        # implicitly return absent-null if non-numeric
    }
    $ox = f($x + NR);
    $oi = f($i);
'
