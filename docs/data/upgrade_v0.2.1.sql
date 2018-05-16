CREATE TABLE pay_trade_chan (
  id            int(10) NOT NULL AUTO_INCREMENT comment '编号',
  trade_no      varchar(45) NOT NULL comment '交易单号',
  pay_chan      int(2) NOT NULL comment '支付途径',
  internal_chan int(1) NOT NULL comment '是否为内置支付途径',
  pay_amount    int(10) NOT NULL comment '支付金额',
  PRIMARY KEY (id)) comment='支付单项';
CREATE TABLE pay_channel (
  id         int(10) NOT NULL AUTO_INCREMENT comment '编号',
  code       varchar(10) NOT NULL comment '支付渠道编码',
  name       int(10) NOT NULL comment '支付渠道名称',
  portal_url varchar(120) NOT NULL comment '支付渠道门户地址',
  PRIMARY KEY (id)) comment='支付通道';
CREATE TABLE pay_order (
  id              int(11) NOT NULL AUTO_INCREMENT comment '编号',
  seller_id       int(11) NOT NULL comment '卖家编号',
  trade_type      varchar(20) NOT NULL comment '交易类型',
  trade_no        varchar(45) NOT NULL comment '交易号',
  subject         varchar(45) NOT NULL comment '支付单详情',
  order_id        int(11) NOT NULL comment '订单号',
  order_type      int(11) NOT NULL comment '支付单的类型，如购物或其他',
  out_order_no    varchar(20) NOT NULL comment '外部订单号',
  buyer_id        int(11) NOT NULL comment '买家编号',
  pay_uid         int(11) NOT NULL comment '支付用户编号',
  item_amount     int(10) NOT NULL comment '商品金额',
  discount_amount int(10) NOT NULL comment '优惠金额 ',
  adjust_amount   int(10) NOT NULL comment '调整金额',
  total_amount    int(10) NOT NULL comment '累计支付金额',
  deduct_amount   int(10) NOT NULL comment '抵扣金额 ',
  procedure_fee   int(10) NOT NULL comment '手续费',
  final_fee       int(10) NOT NULL comment '最终支付金额',
  pay_flag        int(6) NOT NULL comment '可⽤支付方式 ',
  extra_data      varchar(120) NOT NULL comment '其他支付信息',
  trade_channel   int(2) NOT NULL comment '交易支付渠道',
  out_trade_sp    varchar(10) NOT NULL comment '外部交易提供商',
  out_trade_no    varchar(45) NOT NULL comment '外部交易订单号',
  merge_trade_no  varchar(45) NOT NULL comment '合并支付交单单号',
  payment_sign    int(1) NOT NULL comment '可作废',
  state           int(1) NOT NULL comment '订单状态',
  submit_time     int(11) NOT NULL comment '提交时间',
  expires_time    int(11) NOT NULL comment '过期时间',
  paid_time       int(11) NOT NULL comment '支付时间',
  update_time     int(11) NOT NULL comment '更新时间',
  PRIMARY KEY (id)) comment='支付单';

ALTER TABLE `mm_trusted_info`
CHANGE COLUMN `reviewed` `review_state` TINYINT(1) NULL DEFAULT NULL;
ALTER TABLE `mm_levelup`
CHANGE COLUMN `reviewed` `review_state` TINYINT(1) NULL DEFAULT NULL;
ALTER TABLE `mch_enterprise_info`
CHANGE COLUMN `reviewed` `review_state` TINYINT(1) NULL DEFAULT NULL;

