/*
   orders
 */
alter table orders
    alter column accrual type double precision using accrual::double precision;

/*
   withdraw
 */

alter table withdraw
    alter column points type double precision using points::double precision;