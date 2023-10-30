/*
   orders
 */
alter table orders
    alter column accrual type numeric using accrual::numeric;

/*
   withdraw
 */

alter table withdraw
    alter column points type numeric using points::numeric;