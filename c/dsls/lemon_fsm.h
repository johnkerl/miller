#ifndef LEMON_FSM_H
#define LEMON_FSM_H

#include "lemon_structs.h"

// xxx protoize
void FindRulePrecedences(struct lemon*);
void FindFirstSets(struct lemon*);
void FindStates(struct lemon*);
void FindLinks(struct lemon*);
void FindFollowSets(struct lemon*);
void FindActions(struct lemon*);

#endif // LEMON_FSM_H
