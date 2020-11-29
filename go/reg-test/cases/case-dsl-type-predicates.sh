
run_mlr --from $indir/s.dkvp head -n 1 then put -q '
  @is_absent_a       = is_absent       ($a);
  @is_present_a      = is_present      ($a);
  @is_empty_a        = is_empty        ($a);
  @is_not_empty_a    = is_not_empty    ($a);
  @is_null_a         = is_null         ($a);
  @is_not_null_a     = is_not_null     ($a);
  @is_bool_a         = is_bool         ($a);
  @is_boolean_a      = is_boolean      ($a);
  @is_float_a        = is_float        ($a);
  @is_int_a          = is_int          ($a);
  @is_numeric_a      = is_numeric      ($a);
  @is_string_a       = is_string       ($a);
  @is_map_a          = is_map          ($a);
  @is_not_map_a      = is_not_map      ($a);
  @is_empty_map_a    = is_empty_map    ($a);
  @is_nonempty_map_a = is_nonempty_map ($a);
  @is_array_a        = is_array        ($a);
  @is_not_array_a    = is_not_array    ($a);
  dump;
'

run_mlr --from $indir/s.dkvp head -n 1 then put -q '
  @is_absent_x       = is_absent       ($x);
  @is_present_x      = is_present      ($x);
  @is_empty_x        = is_empty        ($x);
  @is_not_empty_x    = is_not_empty    ($x);
  @is_null_x         = is_null         ($x);
  @is_not_null_x     = is_not_null     ($x);
  @is_bool_x         = is_bool         ($x);
  @is_boolean_x      = is_boolean      ($x);
  @is_float_x        = is_float        ($x);
  @is_int_x          = is_int          ($x);
  @is_numeric_x      = is_numeric      ($x);
  @is_string_x       = is_string       ($x);
  @is_map_x          = is_map          ($x);
  @is_not_map_x      = is_not_map      ($x);
  @is_empty_map_x    = is_empty_map    ($x);
  @is_nonempty_map_x = is_nonempty_map ($x);
  @is_array_x        = is_array        ($x);
  @is_not_array_x    = is_not_array    ($x);
  dump;
'

run_mlr --from $indir/s.dkvp head -n 1 then put -q '
  @is_absent_i       = is_absent       ($i);
  @is_present_i      = is_present      ($i);
  @is_empty_i        = is_empty        ($i);
  @is_not_empty_i    = is_not_empty    ($i);
  @is_null_i         = is_null         ($i);
  @is_not_null_i     = is_not_null     ($i);
  @is_bool_i         = is_bool         ($i);
  @is_boolean_i      = is_boolean      ($i);
  @is_float_i        = is_float        ($i);
  @is_int_i          = is_int          ($i);
  @is_numeric_i      = is_numeric      ($i);
  @is_string_i       = is_string       ($i);
  @is_map_i          = is_map          ($i);
  @is_not_map_i      = is_not_map      ($i);
  @is_empty_map_i    = is_empty_map    ($i);
  @is_nonempty_map_i = is_nonempty_map ($i);
  @is_array_i        = is_array        ($i);
  @is_not_array_i    = is_not_array    ($i);
  dump;
'

run_mlr --from $indir/s.dkvp head -n 1 then put -q '
  @is_absent_nonesuch       = is_absent       ($nonesuch);
  @is_present_nonesuch      = is_present      ($nonesuch);
  @is_empty_nonesuch        = is_empty        ($nonesuch);
  @is_not_empty_nonesuch    = is_not_empty    ($nonesuch);
  @is_null_nonesuch         = is_null         ($nonesuch);
  @is_not_null_nonesuch     = is_not_null     ($nonesuch);
  @is_bool_nonesuch         = is_bool         ($nonesuch);
  @is_boolean_nonesuch      = is_boolean      ($nonesuch);
  @is_float_nonesuch        = is_float        ($nonesuch);
  @is_int_nonesuch          = is_int          ($nonesuch);
  @is_numeric_nonesuch      = is_numeric      ($nonesuch);
  @is_string_nonesuch       = is_string       ($nonesuch);
  @is_map_nonesuch          = is_map          ($nonesuch);
  @is_not_map_nonesuch      = is_not_map      ($nonesuch);
  @is_empty_map_nonesuch    = is_empty_map    ($nonesuch);
  @is_nonempty_map_nonesuch = is_nonempty_map ($nonesuch);
  @is_array_nonesuch        = is_array        ($nonesuch);
  @is_not_array_nonesuch    = is_not_array    ($nonesuch);
  dump;
'

run_mlr --from $indir/s.dkvp head -n 1 then put -q '
  @is_absent_dollar_star       = is_absent       ($*);
  @is_present_dollar_star      = is_present      ($*);
  @is_empty_dollar_star        = is_empty        ($*);
  @is_not_empty_dollar_star    = is_not_empty    ($*);
  @is_null_dollar_star         = is_null         ($*);
  @is_not_null_dollar_star     = is_not_null     ($*);
  @is_bool_dollar_star         = is_bool         ($*);
  @is_boolean_dollar_star      = is_boolean      ($*);
  @is_float_dollar_star        = is_float        ($*);
  @is_int_dollar_star          = is_int          ($*);
  @is_numeric_dollar_star      = is_numeric      ($*);
  @is_string_dollar_star       = is_string       ($*);
  @is_map_dollar_star          = is_map          ($*);
  @is_not_map_dollar_star      = is_not_map      ($*);
  @is_empty_map_dollar_star    = is_empty_map    ($*);
  @is_nonempty_map_dollar_star = is_nonempty_map ($*);
  @is_array_dollar_star        = is_array        ($*);
  @is_not_array_dollar_star    = is_not_array    ($*);
  dump;
'

run_mlr --from $indir/s.dkvp head -n 1 then put -q '
  @is_absent_empty_map       = is_absent       ({});
  @is_present_empty_map      = is_present      ({});
  @is_empty_empty_map        = is_empty        ({});
  @is_not_empty_empty_map    = is_not_empty    ({});
  @is_null_empty_map         = is_null         ({});
  @is_not_null_empty_map     = is_not_null     ({});
  @is_bool_empty_map         = is_bool         ({});
  @is_boolean_empty_map      = is_boolean      ({});
  @is_float_empty_map        = is_float        ({});
  @is_int_empty_map          = is_int          ({});
  @is_numeric_empty_map      = is_numeric      ({});
  @is_string_empty_map       = is_string       ({});
  @is_map_empty_map          = is_map          ({});
  @is_not_map_empty_map      = is_not_map      ({});
  @is_empty_map_empty_map    = is_empty_map    ({});
  @is_nonempty_map_empty_map = is_nonempty_map ({});
  @is_array_empty_map        = is_array        ({});
  @is_not_array_empty_map    = is_not_array    ({});
  dump;
'

run_mlr --from $indir/s.dkvp head -n 1 then put -q '
  @is_absent_empty_map       = is_absent       ({});
  @is_present_empty_map      = is_present      ({});
  @is_empty_empty_map        = is_empty        ({});
  @is_not_empty_empty_map    = is_not_empty    ({});
  @is_null_empty_map         = is_null         ({});
  @is_not_null_empty_map     = is_not_null     ({});
  @is_bool_empty_map         = is_bool         ({});
  @is_boolean_empty_map      = is_boolean      ({});
  @is_float_empty_map        = is_float        ({});
  @is_int_empty_map          = is_int          ({});
  @is_numeric_empty_map      = is_numeric      ({});
  @is_string_empty_map       = is_string       ({});
  @is_map_empty_map          = is_map          ({});
  @is_not_map_empty_map      = is_not_map      ({});
  @is_empty_map_empty_map    = is_empty_map    ({});
  @is_nonempty_map_empty_map = is_nonempty_map ({});
  @is_array_empty_map        = is_array        ({});
  @is_not_array_empty_map    = is_not_array    ({});
  dump;
'

run_mlr --from $indir/s.dkvp head -n 1 then put -q '
  $a = "";
  @is_absent_empty       = is_absent       ($a);
  @is_present_empty      = is_present      ($a);
  @is_empty_empty        = is_empty        ($a);
  @is_not_empty_empty    = is_not_empty    ($a);
  @is_null_empty         = is_null         ($a);
  @is_not_null_empty     = is_not_null     ($a);
  @is_bool_empty         = is_bool         ($a);
  @is_boolean_empty      = is_boolean      ($a);
  @is_float_empty        = is_float        ($a);
  @is_int_empty          = is_int          ($a);
  @is_numeric_empty      = is_numeric      ($a);
  @is_string_empty       = is_string       ($a);
  @is_map_empty          = is_map          ($a);
  @is_not_map_empty      = is_not_map      ($a);
  @is_empty_map_empty    = is_empty_map    ($a);
  @is_nonempty_map_empty = is_nonempty_map ($a);
  @is_array_empty        = is_array        ($a);
  @is_not_array_empty    = is_not_array    ($a);
  dump;
'

run_mlr --from $indir/s.dkvp head -n 1 then put -q '
  @is_absent_array       = is_absent       ([1,2,3]);
  @is_present_array      = is_present      ([1,2,3]);
  @is_empty_array        = is_empty        ([1,2,3]);
  @is_not_empty_array    = is_not_empty    ([1,2,3]);
  @is_null_array         = is_null         ([1,2,3]);
  @is_not_null_array     = is_not_null     ([1,2,3]);
  @is_bool_array         = is_bool         ([1,2,3]);
  @is_boolean_array      = is_boolean      ([1,2,3]);
  @is_float_array        = is_float        ([1,2,3]);
  @is_int_array          = is_int          ([1,2,3]);
  @is_numeric_array      = is_numeric      ([1,2,3]);
  @is_string_array       = is_string       ([1,2,3]);
  @is_map_array          = is_map          ([1,2,3]);
  @is_not_map_array      = is_not_map      ([1,2,3]);
  @is_empty_map_array    = is_empty_map    ([1,2,3]);
  @is_nonempty_map_array = is_nonempty_map ([1,2,3]);
  @is_array_array        = is_array        ([1,2,3]);
  @is_not_array_array    = is_not_array    ([1,2,3]);
  dump;
'

run_mlr --from $indir/s.dkvp head -n 1 then put -q '
  v = [1,2,3];
  @is_absent_array_in_bounds       = is_absent       (v[1]);
  @is_present_array_in_bounds      = is_present      (v[1]);
  @is_empty_array_in_bounds        = is_empty        (v[1]);
  @is_not_empty_array_in_bounds    = is_not_empty    (v[1]);
  @is_null_array_in_bounds         = is_null         (v[1]);
  @is_not_null_array_in_bounds     = is_not_null     (v[1]);
  @is_bool_array_in_bounds         = is_bool         (v[1]);
  @is_boolean_array_in_bounds      = is_boolean      (v[1]);
  @is_float_array_in_bounds        = is_float        (v[1]);
  @is_int_array_in_bounds          = is_int          (v[1]);
  @is_numeric_array_in_bounds      = is_numeric      (v[1]);
  @is_string_array_in_bounds       = is_string       (v[1]);
  @is_map_array_in_bounds          = is_map          (v[1]);
  @is_not_map_array_in_bounds      = is_not_map      (v[1]);
  @is_empty_map_array_in_bounds    = is_empty_map    (v[1]);
  @is_nonempty_map_array_in_bounds = is_nonempty_map (v[1]);
  @is_array_in_bounds_array        = is_array        (v[1]);
  @is_not_array_in_bounds_array    = is_not_array    (v[1]);
  dump;
'

run_mlr --from $indir/s.dkvp head -n 1 then put -q '
  v = [1,2,3];
  @is_absent_array_out_of_bounds       = is_absent       (v[4]);
  @is_present_array_out_of_bounds      = is_present      (v[4]);
  @is_empty_array_out_of_bounds        = is_empty        (v[4]);
  @is_not_empty_array_out_of_bounds    = is_not_empty    (v[4]);
  @is_null_array_out_of_bounds         = is_null         (v[4]);
  @is_not_null_array_out_of_bounds     = is_not_null     (v[4]);
  @is_bool_array_out_of_bounds         = is_bool         (v[4]);
  @is_boolean_array_out_of_bounds      = is_boolean      (v[4]);
  @is_float_array_out_of_bounds        = is_float        (v[4]);
  @is_int_array_out_of_bounds          = is_int          (v[4]);
  @is_numeric_array_out_of_bounds      = is_numeric      (v[4]);
  @is_string_array_out_of_bounds       = is_string       (v[4]);
  @is_map_array_out_of_bounds          = is_map          (v[4]);
  @is_not_map_array_out_of_bounds      = is_not_map      (v[4]);
  @is_empty_map_array_out_of_bounds    = is_empty_map    (v[4]);
  @is_nonempty_map_array_out_of_bounds = is_nonempty_map (v[4]);
  @is_array_out_of_bounds_array        = is_array        (v[4]);
  @is_not_array_out_of_bounds_array    = is_not_array    (v[4]);
  dump;
'

run_mlr --from $indir/s.dkvp head -n 1 then put -q '
  v = [1,2,3];
  @is_absent_array_out_of_bounds       = is_absent       (v[0]);
  @is_present_array_out_of_bounds      = is_present      (v[0]);
  @is_empty_array_out_of_bounds        = is_empty        (v[0]);
  @is_not_empty_array_out_of_bounds    = is_not_empty    (v[0]);
  @is_null_array_out_of_bounds         = is_null         (v[0]);
  @is_not_null_array_out_of_bounds     = is_not_null     (v[0]);
  @is_bool_array_out_of_bounds         = is_bool         (v[0]);
  @is_boolean_array_out_of_bounds      = is_boolean      (v[0]);
  @is_float_array_out_of_bounds        = is_float        (v[0]);
  @is_int_array_out_of_bounds          = is_int          (v[0]);
  @is_numeric_array_out_of_bounds      = is_numeric      (v[0]);
  @is_string_array_out_of_bounds       = is_string       (v[0]);
  @is_map_array_out_of_bounds          = is_map          (v[0]);
  @is_not_map_array_out_of_bounds      = is_not_map      (v[0]);
  @is_empty_map_array_out_of_bounds    = is_empty_map    (v[0]);
  @is_nonempty_map_array_out_of_bounds = is_nonempty_map (v[0]);
  @is_array_out_of_bounds_array        = is_array        (v[0]);
  @is_not_array_out_of_bounds_array    = is_not_array    (v[0]);
  dump;
'

run_mlr --from $indir/s.dkvp head -n 1 then put -q '
  @asserting_present_x                = asserting_present      ($x);
  @asserting_not_empty_x              = asserting_not_empty    ($x);
  @asserting_not_null_x               = asserting_not_null     ($x);
  @asserting_float_x                  = asserting_float        ($x);
  @asserting_int_i                    = asserting_int          ($i);
  @asserting_numeric_x                = asserting_numeric      ($x);
  @asserting_numeric_i                = asserting_numeric      ($i);
  @asserting_string_b                 = asserting_string       ($b);
  @asserting_map_dollar_star          = asserting_map          ($*);
  @asserting_not_map_x                = asserting_not_map      ($x);
  @asserting_empty_map_curlies        = asserting_empty_map    ({});
  @asserting_nonempty_map_dollar_star = asserting_nonempty_map ($*);
  @asserting_array_braces             = asserting_array        ([]);
  @asserting_not_array_x              = asserting_not_array    ($x);
  dump;
'

mlr_expect_fail --from $indir/s.dkvp head -n 1 then put -q '
  @asserting_absent_x = asserting_absent($x);
'


mlr_expect_fail --from $indir/s.dkvp head -n 1 then put -q '
  @asserting_empty_x = asserting_empty($x);
  dump;
'
