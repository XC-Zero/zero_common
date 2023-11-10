create table sales_order_source
(
    abbrevtype                            STRING,
    actualshipdate                        STRING,
    altsalestotal                         STRING,
    approvalstatus                        STRING,
    billingstatus                         STRING,
    closedate                             timestamp,
    createdby                             bigint,
    createddate                           timestamp,
    currency                              bigint,
    custbody_aio_account                  bigint,
    custbody_aio_marketplaceid            bigint,
    custbody_aio_payment_time             timestamp,
    custbody_aio_report_settlement_date   STRING,
    custbody_aio_s_fulfillment_channel    STRING,
    custbody_aio_s_is_business_order      STRING,
    custbody_aio_s_is_premium_order       STRING,
    custbody_aio_s_is_prime               STRING,
    custbody_aio_s_l_u_date               STRING,
    custbody_aio_s_order_status           STRING,
    custbody_aio_s_order_type             STRING,
    custbody_aio_s_p_date                 STRING,
    custbody_aio_s_sales_channel          STRING,
    custbody_aio_settlement_id            STRING,
    custbody_aio_shipping_category        STRING,
    custbody_aio_sync_multi_channel       STRING,
    custbody_aio_sync_to_platform         STRING,
    custbody_delivery_channel             STRING,
    custbody_dps_quality_inspection_order STRING,
    custbody_is_back_order                STRING,
    custbody_type_of_shipping             STRING,
    custbody_wave_create_flag             STRING,
    custbodycustomreport_site             STRING,
    custbodycustomrole_if                 STRING,
    entity                                bigint,
    foreigntotal                          DECIMAL,
    id                                    bigint,
    intercostatus                         STRING,
    lastmodifiedby                        bigint,
    lastmodifieddate                      timestamp,
    memo                                  STRING,
    nextapprover                          STRING,
    ordpicked                             STRING,
    otherrefnum                           STRING,
    paymentmethod                         bigint,
    recordtype                            STRING,
    shipcarrier                           STRING,
    shipcomplete                          STRING,
    shipdate                              timestamp,
    shippingaddress                       bigint,
    source                                STRING,
    status                                STRING,
    totalcostestimate                     DECIMAL,
    trackingnumberlist                    STRING,
    trandate                              timestamp,
    tranid                                STRING,
    transactionnumber                     STRING,
    type                                  STRING,
    uniquekey                             bigint,
    linelastmodifieddate                  timestamp,
    subsidiary                            bigint,
    rateamount                            DECIMAL,
    debitforeignamount                    DECIMAL,
    linesequencenumber                    bigint,
    accountinglinetype                    STRING,
    commitinventory                       DECIMAL,
    commitmentfirm                        STRING,
    custcol_aio_accepted_quantity         STRING,
    custcol_aio_already_taxed             STRING,
    custcol_aio_amazon_asin               STRING,
    custcol_aio_amazon_msku               STRING,
    custcol_aio_order_item_id             STRING,
    custcol_aio_s_gift_wrap_price         STRING,
    custcol_aio_s_gift_wrap_tax           DECIMAL,
    custcol_aio_s_item_tax                DECIMAL,
    custcol_aio_s_item_title              STRING,
    custcol_aio_s_points_amount           DECIMAL,
    custcol_aio_s_promotion_discount      DECIMAL,
    custcol_aio_s_promotion_discount_tax  DECIMAL,
    custcol_aio_s_promotion_ids           STRING,
    custcol_aio_s_shipping_discount       DECIMAL,
    custcol_aio_s_shipping_discount_tax   DECIMAL,
    custcol_aio_s_shipping_price          DECIMAL,
    custcol_aio_s_shipping_tax            DECIMAL,
    custcol_aio_ship_type                 STRING,
    custcol_expected_receipt_date         STRING,
    custcol_is_expedited                  STRING,
    custcol_logistics_method              STRING,
    custcol_model_order                   STRING,
    department                            DECIMAL,
    estgrossprofit                        DECIMAL,
    expectedreceiptdate                   STRING,
    expenseaccount                        DECIMAL,
    hasfulfillableitems                   STRING,
    isinventoryaffecting                  STRING,
    item                                  DECIMAL,
    itemtype                              STRING,
    location                              DECIMAL,
    mainline                              STRING,
    matchbilltoreceipt                    STRING,
    netamount                             DECIMAL,
    oldcommitmentfirm                     STRING,
    orderpriority                         STRING,
    price                                 DECIMAL,
    processedbyrevcommit                  STRING,
    quantity                              DECIMAL,
    quantitybackordered                   STRING,
    quantitybilled                        DECIMAL,
    quantitypacked                        DECIMAL,
    quantitypicked                        DECIMAL,
    quantityrejected                      DECIMAL,
    quantityshiprecv                      DECIMAL,
    transactionlinetype                   STRING,
    transferorderitemlineid               STRING
)
WITH ('connector' = 'mysql-cdc',
    'hostname' = 'tessan-mysql.mysql.polardb.rds.aliyuncs.com',
    'port' = '3306',
    'username' = 'tessan_erp_all',
    'password' = 'f6Yc27Ds@a3GeHN',
    'database-name' = 'tessan_erp_show_data',
    'table-name' = 'tessan_sales_order',
    'server-time-zone' = 'UTC',
    'scan.incremental.snapshot.enabled'='false');



create table sales_order_target
(
    abbrevtype                            STRING,
    actualshipdate                        STRING,
    altsalestotal                         STRING,
    approvalstatus                        STRING,
    billingstatus                         STRING,
    closedate                             timestamp,
    createdby                             bigint,
    createddate                           timestamp,
    currency                              bigint,
    custbody_aio_account                  bigint,
    custbody_aio_marketplaceid            bigint,
    custbody_aio_payment_time             timestamp,
    custbody_aio_report_settlement_date   STRING,
    custbody_aio_s_fulfillment_channel    STRING,
    custbody_aio_s_is_business_order      STRING,
    custbody_aio_s_is_premium_order       STRING,
    custbody_aio_s_is_prime               STRING,
    custbody_aio_s_l_u_date               STRING,
    custbody_aio_s_order_status           STRING,
    custbody_aio_s_order_type             STRING,
    custbody_aio_s_p_date                 STRING,
    custbody_aio_s_sales_channel          STRING,
    custbody_aio_settlement_id            STRING,
    custbody_aio_shipping_category        STRING,
    custbody_aio_sync_multi_channel       STRING,
    custbody_aio_sync_to_platform         STRING,
    custbody_delivery_channel             STRING,
    custbody_dps_quality_inspection_order STRING,
    custbody_is_back_order                STRING,
    custbody_type_of_shipping             STRING,
    custbody_wave_create_flag             STRING,
    custbodycustomreport_site             STRING,
    custbodycustomrole_if                 STRING,
    entity                                bigint,
    foreigntotal                          DECIMAL,
    id                                    bigint,
    intercostatus                         STRING,
    lastmodifiedby                        bigint,
    lastmodifieddate                      timestamp,
    memo                                  STRING,
    nextapprover                          STRING,
    ordpicked                             STRING,
    otherrefnum                           STRING,
    paymentmethod                         bigint,
    recordtype                            STRING,
    shipcarrier                           STRING,
    shipcomplete                          STRING,
    shipdate                              timestamp,
    shippingaddress                       bigint,
    source                                STRING,
    status                                STRING,
    totalcostestimate                     DECIMAL,
    trackingnumberlist                    STRING,
    trandate                              timestamp,
    tranid                                STRING,
    transactionnumber                     STRING,
    type                                  STRING,
    uniquekey                             bigint,
    linelastmodifieddate                  timestamp,
    subsidiary                            bigint,
    rateamount                            DECIMAL,
    debitforeignamount                    DECIMAL,
    linesequencenumber                    bigint,
    accountinglinetype                    STRING,
    commitinventory                       DECIMAL,
    commitmentfirm                        STRING,
    custcol_aio_accepted_quantity         STRING,
    custcol_aio_already_taxed             STRING,
    custcol_aio_amazon_asin               STRING,
    custcol_aio_amazon_msku               STRING,
    custcol_aio_order_item_id             STRING,
    custcol_aio_s_gift_wrap_price         STRING,
    custcol_aio_s_gift_wrap_tax           DECIMAL,
    custcol_aio_s_item_tax                DECIMAL,
    custcol_aio_s_item_title              STRING,
    custcol_aio_s_points_amount           DECIMAL,
    custcol_aio_s_promotion_discount      DECIMAL,
    custcol_aio_s_promotion_discount_tax  DECIMAL,
    custcol_aio_s_promotion_ids           STRING,
    custcol_aio_s_shipping_discount       DECIMAL,
    custcol_aio_s_shipping_discount_tax   DECIMAL,
    custcol_aio_s_shipping_price          DECIMAL,
    custcol_aio_s_shipping_tax            DECIMAL,
    custcol_aio_ship_type                 STRING,
    custcol_expected_receipt_date         STRING,
    custcol_is_expedited                  STRING,
    custcol_logistics_method              STRING,
    custcol_model_order                   STRING,
    department                            DECIMAL,
    estgrossprofit                        DECIMAL,
    expectedreceiptdate                   STRING,
    expenseaccount                        DECIMAL,
    hasfulfillableitems                   STRING,
    isinventoryaffecting                  STRING,
    item                                  DECIMAL,
    itemtype                              STRING,
    location                              DECIMAL,
    mainline                              STRING,
    matchbilltoreceipt                    STRING,
    netamount                             DECIMAL,
    oldcommitmentfirm                     STRING,
    orderpriority                         STRING,
    price                                 DECIMAL,
    processedbyrevcommit                  STRING,
    quantity                              DECIMAL,
    quantitybackordered                   STRING,
    quantitybilled                        DECIMAL,
    quantitypacked                        DECIMAL,
    quantitypicked                        DECIMAL,
    quantityrejected                      DECIMAL,
    quantityshiprecv                      DECIMAL,
    transactionlinetype                   STRING,
    transferorderitemlineid               STRING
)
WITH ('connector' = 'mysql-cdc',
    'hostname' = '192.168.12.225',
    'port' = '3306',
    'username' = 'root',
    'password' = 'root123456',
    'database-name' = 'temp',
    'table-name' = 'tessan_sales_order',
    'server-time-zone' = 'UTC',
    'scan.incremental.snapshot.enabled'='false');

insert into sales_order_target
select *
from sales_order_source
where custbody_aio_payment_time >= '2023-04-01'
  and custbody_aio_payment_time <= '2023-04-30';
