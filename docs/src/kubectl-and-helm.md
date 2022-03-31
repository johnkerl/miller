<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-main-flag-list/index.html">Flags</a>
&nbsp;
<a class="quicklink" href="../reference-verbs/index.html">Verbs</a>
&nbsp;
<a class="quicklink" href="../reference-dsl-builtin-functions/index.html">Functions</a>
&nbsp;
<a class="quicklink" href="../glossary/index.html">Glossary</a>
&nbsp;
<a class="quicklink" href="../release-docs/index.html">Release docs</a>
</span>
</div>
# Processing Kubectl and Helm output

The `kubectl` and `helm` commands produce tabular-looking output, which Miller can parse -- however,
there's a bit of whitespace-handling to be dealt with first.

## Whitespace structure of kubectl

The output of `kubectl` looks tabular -- so [PPRINT format](file-formats.md#pprint-pretty-printed-tabular) is perhaps a good choice.

<pre class="pre-non-highlight-non-pair">
$ kubectl -n my-namespace get pods | head
NAME                      READY STATUS    RESTARTS AGE
app-5mjwm4-274754k8468    0/1   Completed 0        6m51s
app-5mjwm4-274754vdfnf    0/1   Completed 0        6m50s
app-5mjwm4-0              1/1   Running   0        6h8m
app-5mjwm4-27475nt9cc     0/1   Completed 0        6m53s
app-5mjwm4-27474454-dc7wq 0/1   Error     0        16h
app-5mjwm4-27475416-tv2ff 0/1   Completed 0        56s
app-5mjwm4-2747541t7lgk   0/1   Completed 0        115s
app-5mjwm4-27475245-7sg9r 0/1   Completed 0        171m
app-5mjwm4-27475410-k4gcr 0/1   Completed 0        6m52s
</pre>

We can verify this using `kubectl -n my-namespace get pods | vim -`, then `:set list` within `vim`;
or, perhaps piping the output to `cat -t`, or [bat -A](https://github.com/sharkdp/bat) -- in any
case what looks like whitespace is really all space characters.

To double-check, it's helpful to run tabular-looking output through a format-converter, and make
sure the column headers are being correctly identified as keys, and the remaining lines are being
correctly identified as values:

<pre class="pre-non-highlight-non-pair">
$ kubectl -n my-namespace get pods | mlr --ipprint --ojson head -n 1
[
{
  "NAME": "app-5mjwm4-274754k8468",
  "READY": "0/1",
  "STATUS": "Completed",
  "RESTARTS": 0,
  "AGE": "14m"
}
]
</pre>

## Sorting/filtering

Suppose we want to sort the information for non-completed pods by age. We can use
[dhms2sec](reference-dsl-builtin-functions.md#dhms2sec) to turn the `AGE` into something sortable.

<pre class="pre-non-highlight-non-pair">
$ kubectl -n service-xyz get pods \
  | mlr --pprint \
    filter '$STATUS != "Completed"' \
    then put '$AGESEC = dhms2sec($AGE)' \
    then sort -n AGESEC
NAME                              READY STATUS  RESTARTS AGE   AGESEC
app1-1500-5mjwm4-0                1/1   Running 0        6h22m 22920
app1-1624-6dh711-0                1/1   Running 0        6h27m 23220
app1-1500-pqb9b4-0                1/1   Running 0        6h30m 23400
app1-gbwuwi-2747495lbtzg          0/1   Error   0        7h59m 28740
app1-gbwuwi-0                     1/1   Running 0        8h    28800
app1-gbwuwi-27474955r8gq          0/1   Error   0        8h    28800
app1-gbwuwi-27474956rps8          0/1   Error   0        8h    28800
app1-gbwuwi-2747495q7fnz          0/1   Error   0        8h    28800
app1-gbwuwi-2747495vnxgn          0/1   Error   0        8h    28800
app1-gbwuwi-674ddcfd89-2jt64      2/2   Running 0        8h    28800
app3-5c79574b69-8njgr             2/2   Running 0        9h    32400
app3-5c79574b69-np2qj             2/2   Running 0        9h    32400
app3-a56i7c-0                     1/1   Running 0        13h   46800
app3-a56i7c-587dfc99cf-zrr4t      2/2   Running 0        13h   46800
app2-1500-pqb9b4-274746pfbfd      0/1   Error   0        13h   46800
app2-1500-pqb9b4-274746jtz8t      0/1   Error   0        13h   46800
app2-1500-pqb9b4-274746pmmhq      0/1   Error   0        13h   46800
app2-1500-pqb9b4-27474624h8fp     0/1   Error   0        13h   46800
app2-1500-pqb9b4-2747462d8n96     0/1   Error   0        13h   46800
app2-1500-pqb9b4-2747462xnmcf     0/1   Error   0        13h   46800
app2-1500-pqb9b4-27474630-95668   0/1   Error   0        13h   46800
app1-1500-pqb9b4-sr5vd            2/2   Running 0        13h   46800
app1-1500-5mjwm4-27474454-dc7wq   0/1   Error   0        16h   57600
app1-1500-5mjwm4-667c6fc66d-b97m9 2/2   Running 0        16h   57600
app1-1624-6dh711-2747435h42j      0/1   Error   0        17h   61200
app1-1624-6dh711-27474370-ph25r   0/1   Error   0        17h   61200
app1-1624-6dh711-74fb5cf9d6-cl5tq 2/2   Running 0        17h   61200
</pre>

## Whitespace structure of helm list

The output of `helm list` is a bit fussier. Here it's already clear that something's amiss, since not everything lines up:

<pre class="pre-non-highlight-non-pair">
$ helm list
NAME                      NAMESPACE   REVISION  UPDATED                                 STATUS    CHART                     APP VERSION
appdev-an-sc-1500-5mjwm4  service-xyz 1         2022-03-28 11:33:05.389975262 +0000 UTC deployed  appdev-cloud-test-7.1.12
appdev-exyzv-load-a56i7c  service-xyz 1         2022-03-28 14:45:35.44317196 +0000 UTC  deployed  appdev-cloud-test-7.1.12
appdev-sa-sc-1500-pqb9b4  service-xyz 1         2022-03-28 14:24:33.978580048 +0000 UTC deployed  appdev-cloud-test-7.1.12
appdev-sa-sc-1624-6dh711  service-xyz 1         2022-03-28 10:09:05.966332699 +0000 UTC deployed  appdev-cloud-test-7.1.12
appdev-wertzxyffa-gbwuwi  service-xyz 1         2022-03-28 19:47:34.96763583 +0000 UTC  deployed  appdev-cloud-test-7.1.12
staging                   service-xyz 797       2022-03-28 18:39:34.005120936 +0000 UTC deployed  appdev-cloud-test-7.1.12
</pre>

This isn't likely to be PPRINT format, as we soon see. Also note that the space before `+0000` is an issue.

<pre class="pre-non-highlight-non-pair">
$ helm list | mlr --ipprint --ojson cat
mlr :  mlr: CSV header/data length mismatch 7 != 5 at filename (stdin) line  2.
</pre>

Running through `bat -A` or `cat -t` shows an issue. Namely, the Helm authors are mixing tabs and spaces -- `cat -t` shows tabs as `^I`:

<pre class="pre-non-highlight-non-pair">
$ helm list | cat -t
NAME                    ^INAMESPACE  ^IREVISION^IUPDATED                                ^ISTATUS  ^ICHART                   ^IAPP VERSION
appdev-an-sc-1500-5mjwm4^Iservice-xyz^I1       ^I2022-03-28 11:33:05.389975262 +0000 UTC^Ideployed^Iappdev-cloud-test-7.1.12^I
appdev-exyzv-load-a56i7c^Iservice-xyz^I1       ^I2022-03-28 14:45:35.44317196 +0000 UTC ^Ideployed^Iappdev-cloud-test-7.1.12^I
appdev-sa-sc-1500-pqb9b4^Iservice-xyz^I1       ^I2022-03-28 14:24:33.978580048 +0000 UTC^Ideployed^Iappdev-cloud-test-7.1.12^I
appdev-sa-sc-1624-6dh711^Iservice-xyz^I1       ^I2022-03-28 10:09:05.966332699 +0000 UTC^Ideployed^Iappdev-cloud-test-7.1.12^I
appdev-wertzxyffa-gbwuwi^Iservice-xyz^I1       ^I2022-03-28 19:47:34.96763583 +0000 UTC ^Ideployed^Iappdev-cloud-test-7.1.12^I
staging                 ^Iservice-xyz^I797     ^I2022-03-28 18:39:34.005120936 +0000 UTC^Ideployed^Iappdev-cloud-test-7.1.12^I
</pre>

This mix of tabs and spaces, while not PPRINT, also isn't quite TSV either. As above, it's helpful to run tabular-looking data through a format-converter
to see how it's structured:

<pre class="pre-non-highlight-non-pair">
$ helm list | mlr --itsv --ojson head -n 1
[
{
  "NAME                    ": "appdev-an-sc-1500-5mjwm4",
  "NAMESPACE  ": "service-xyz",
  "REVISION": "1       ",
  "UPDATED                                ": "2022-03-28 11:33:05.389975262 +0000 UTC",
  "STATUS  ": "deployed",
  "CHART                   ": "appdev-cloud-test-7.1.12",
  "APP VERSION": "           "
}
]
</pre>

A solution here is Miller's 
[clean-whitespace verb](reference-verbs.md#clean-whitespace):

<pre class="pre-non-highlight-non-pair">
$ helm list | mlr --itsv --ojson clean-whitespace then head -n 1
[
{
  "NAME": "appdev-an-sc-1500-5mjwm4",
  "NAMESPACE": "service-xyz",
  "REVISION": "1       ",
  "UPDATED": "2022-03-28 11:33:05.389975262 +0000 UTC",
  "STATUS ": "deployed",
  "CHART": "appdev-cloud-test-7.1.12",
  "APP VERSION": ""
}
]
</pre>

Now we have the keys and values correctly identified within the tabular-looking data.

## Sorting/filtering

To find oldest items, it would suffice to sort by the `UPDATED` column, as that sorts lexically.
However, let's parse the timestamps and compute their ages from the present:

<pre class="pre-non-highlight-non-pair">
$ helm list \
  | mlr --itsv --opprint clean-whitespace \
    then put '$AGESEC = int(systime() - strptime($UPDATED, "%Y-%m-%d %H:%M:%S.%f +0000 UTC"))' \
    then sort -n AGESEC \
    then cut -x -f 'APP VERSION,UPDATED'
NAME                     NAMESPACE   REVISION STATUS   CHART                    AGESEC
appdev-sa-sc-1624-6dh711 service-xyz 1        deployed appdev-cloud-test-7.1.12 30874
appdev-an-sc-1500-5mjwm4 service-xyz 797      deployed appdev-cloud-test-7.1.12 34955
appdev-sa-sc-1500-pqb9b4 service-xyz 1        deployed appdev-cloud-test-7.1.12 48993
appdev-xxyzv-load-a56i7c service-xyz 1        deployed appdev-cloud-test-7.1.12 50255
staging                  service-xyz 1        deployed appdev-cloud-test-7.1.12 60543
appdev-wertzxyffa-gbwuwi service-xyz 1        deployed appdev-cloud-test-7.1.12 65583
</pre>

## Extracting fields to be acted on

Switching to [NIDX format](file-formats.md#nidx-index-numbered-toolkit-style) lets us extract fields and pass them onto other commands -- e.g. `helm uninstall`.
We just need to switch the output format to `--onidx`, then cut out the `NAME` field. (Maybe add `then filter $AGESEC > 86400` or somesuch.)

<pre class="pre-non-highlight-non-pair">
$ helm list \
  | mlr --itsv --onidx clean-whitespace \
    then put '$UPDATED = ssub($UPDATED, " +0000 UTC", "")' \
    then put '$AGESEC = int(systime() - strptime($UPDATED, "%Y-%m-%d %H:%M:%S.%f"))' \
    then sort -n AGESEC \
    then cut -f NAME \
  | tee names.txt
appdev-sa-sc-1624-6dh711
appdev-an-sc-1500-5mjwm4
appdev-sa-sc-1500-pqb9b4
appdev-xxyzv-load-a56i7c
staging
appdev-wertzxyffa-gbwuwi
</pre>

Then

<pre class="pre-non-highlight-non-pair">
$ for name in $(cat names.txt); do helm uninstall $name; done
</pre>

