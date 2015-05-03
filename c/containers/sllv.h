// Singly-linked list of void-star, with tail for append.
#ifndef SLLV_H
#define SLLV_H

typedef struct _sllve_t {
	void* pvdata;
	struct _sllve_t *pnext;
} sllve_t;

typedef struct _sllv_t {
	sllve_t *phead;
	sllve_t *ptail;
	int length;
} sllv_t;

sllv_t* sllv_alloc();
void    sllv_free(sllv_t* plist);
sllv_t* sllv_single(void* pvdata);
void    sllv_add(sllv_t* plist, void* pvdata);
void    sllv_reverse(sllv_t* plist);

void*   sllv_pop(sllv_t* plist);

// xxx cmt mem-mgt
sllv_t* sllv_append(sllv_t* pa, sllv_t* pb);

#endif // SLLV_H

