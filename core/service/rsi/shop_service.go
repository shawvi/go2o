/**
 * Copyright 2015 @ to2.net.
 * name : content_service
 * author : jarryliu
 * date : -- :
 * description :
 * history :
 */
package rsi

import (
	"context"
	"errors"
	"go2o/core/domain/interface/merchant"
	"go2o/core/domain/interface/merchant/shop"
	"go2o/core/dto"
	"go2o/core/infrastructure/format"
	"go2o/core/query"
	"go2o/core/service/auto_gen/rpc/shop_service"
	"go2o/core/service/auto_gen/rpc/ttype"
	"go2o/core/service/thrift/parser"
	"go2o/core/variable"
)

var _ shop_service.ShopService = new(shopServiceImpl)

type shopServiceImpl struct {
	repo     shop.IShopRepo
	mchRepo  merchant.IMerchantRepo
	shopRepo shop.IShopRepo
	query    *query.ShopQuery
	serviceUtil
}

func (si *shopServiceImpl) QueryShopIdByHost(ctx context.Context, host string) (r int32, err error) {
	_, shopId := si.query.QueryShopIdByHost(host)
	return shopId, nil
}

func (si *shopServiceImpl) GetShop(ctx context.Context, shopId int32) (r *shop_service.SShop, err error) {
	sp := si.shopRepo.GetOnlineShop(int(shopId))
	if sp != nil {
		return si.parseShop(sp), nil
	}
	return nil, nil
}

func (si *shopServiceImpl) GetVendorShop(ctx context.Context, vendorId int32) (r *shop_service.SShop, err error) {
	sp := si.shopRepo.GetOnlineShopOfMerchant(int(vendorId))
	if sp != nil {
		return si.parseShop(sp), nil
	}
	return nil, nil
}

// 根据主机头获取店铺编号
func (si *shopServiceImpl) QueryShopByHost(ctx context.Context, host string) (r int32, err error) {
	_, shopId := si.query.QueryShopIdByHost(host)
	return shopId, nil
}

func NewShopService(rep shop.IShopRepo, mchRepo merchant.IMerchantRepo,
	shopRepo shop.IShopRepo, query *query.ShopQuery) *shopServiceImpl {
	return &shopServiceImpl{
		repo:     rep,
		mchRepo:  mchRepo,
		shopRepo: shopRepo,
		query:    query,
	}
}

// 获取门店
func (si *shopServiceImpl) GetStore(ctx context.Context, storeId int32) (*shop_service.SStore, error) {
	panic("返回门店")
	//mch := si.mchRepo.GetMerchant(int(storeId))
	//if mch != nil {
	//	shop := mch.ShopManager().GetOnlineShop()
	//	if shop != nil {
	//		return parser.ParseOnlineShop(shop), nil
	//	}
	//}
	//return nil, nil
}

func (si *shopServiceImpl) GetStoreById(ctx context.Context, shopId int32) (*shop_service.SStore, error) {
	vendorId := si.query.GetMerchantId(shopId)
	return si.GetStore(ctx, vendorId)
}

// 打开或关闭商店
func (si *shopServiceImpl) TurnShop(ctx context.Context, shopId int32, on bool, reason string) (*ttype.Result_, error) {
	var err error
	sp := si.repo.GetShop(int(shopId))
	if sp == nil {
		err = shop.ErrNoSuchShop
	} else {
		if on {
			err = sp.TurnOn()
		} else {
			err = sp.TurnOff(reason)
		}
	}
	return si.result(err), nil
}

// 设置商店是否营业
func (si *shopServiceImpl) OpenShop(ctx context.Context, shopId int32, on bool, reason string) (*ttype.Result_, error) {
	var err error
	sp := si.repo.GetShop(int(shopId))
	if sp == nil {
		err = shop.ErrNoSuchShop
	} else {
		if on {
			err = sp.Opening()
		} else {
			err = sp.Pause()
		}
	}
	return si.result(err), nil
}

func (si *shopServiceImpl) getMerchantId(shopId int32) int32 {
	return si.query.GetMerchantId(shopId)
}

func (si *shopServiceImpl) GetMerchantId(shopId int32) int32 {
	return si.query.GetMerchantId(shopId)
}

// 获取商店的数据
func (si *shopServiceImpl) GetShopData(mchId, shopId int32) *shop.ComplexShop {
	mch := si.mchRepo.GetMerchant(int(mchId))
	sp := mch.ShopManager().GetShop(int(shopId))
	if sp != nil {
		return sp.Data()
	}
	return nil
}

func (si *shopServiceImpl) GetShopValueById(mchId, shopId int32) *shop.Shop {
	mch := si.mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		v := mch.ShopManager().GetShop(int(shopId)).GetValue()
		return &v
	}
	return nil
}

// 保存线上商店
func (si *shopServiceImpl) SaveStore(s *shop_service.SStore) error {
	mch := si.mchRepo.GetMerchant(int(s.VendorId))
	if mch != nil {
		v, v1 := parser.Parse2OnlineShop(s)
		mgr := mch.ShopManager()
		sp := mgr.GetOnlineShop()
		// 创建商店
		if sp == nil {
			sp = mgr.CreateShop(v)
		}
		err := sp.SetValue(v)
		if err == nil {
			ofs := sp.(shop.IOnlineShop)
			err = ofs.SetShopValue(v1)
			if err == nil {
				err = sp.Save()
			}
		}
		return err
	}
	return merchant.ErrNoSuchMerchant
}

// 保存门店
func (si *shopServiceImpl) SaveOfflineShop(s *shop.Shop, v *shop.OfflineShop) error {
	mch := si.mchRepo.GetMerchant(int(s.VendorId))
	if mch != nil {
		mgr := mch.ShopManager()
		var sp shop.IShop
		if s.Id > 0 {
			// 保存商店
			sp = mgr.GetShop(int(s.Id))
		} else {
			//创建商店
			sp = mgr.CreateShop(s)
		}
		err := sp.SetValue(s)
		if err == nil {
			ofs := sp.(shop.IOfflineShop)
			err = ofs.SetShopValue(v)
			if err == nil {
				err = sp.Save()
			}
		}
		return err
	}
	return merchant.ErrNoSuchMerchant
}

func (si *shopServiceImpl) SaveShop(mchId int32, v *shop.Shop) (int32, error) {
	mch := si.mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		var shop shop.IShop
		if v.Id > 0 {
			shop = mch.ShopManager().GetShop(int(v.Id))
			if shop == nil {
				return 0, errors.New("门店不存在")
			}
		} else {
			shop = mch.ShopManager().CreateShop(v)
		}
		err := shop.SetValue(v)
		if err != nil {
			return v.Id, err
		}
		err = shop.Save()
		return int32(shop.GetDomainId()), err
	}
	return 0, merchant.ErrNoSuchMerchant
}

func (si *shopServiceImpl) DeleteShop(mchId, shopId int32) error {
	mch := si.mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		return mch.ShopManager().DeleteShop(shopId)
	}
	return merchant.ErrNoSuchMerchant
}

// 获取线上商城配置
func (si *shopServiceImpl) GetOnlineShopConf(shopId int32) *shop.OnlineShop {
	mchId := si.getMerchantId(shopId)
	mch := si.mchRepo.GetMerchant(int(mchId))
	if mch != nil {
		s := mch.ShopManager().GetShop(int(shopId))
		if s == nil {
			v := s.(shop.IOnlineShop).GetShopValue()
			return &v
		}
	}
	return nil
}

// 获取商城
func (si *shopServiceImpl) GetOnlineShops(vendorId int32) []*shop.Shop {
	mch := si.mchRepo.GetMerchant(int(vendorId))
	shops := mch.ShopManager().GetShops()
	sv := make([]*shop.Shop, 0)
	for _, v := range shops {
		if v.Type() == shop.TypeOnlineShop {
			vv := v.GetValue()
			sv = append(sv, &vv)
		}
	}
	return sv
}

// 获取指定的营业中的店铺
func (si *shopServiceImpl) PagedOnBusinessOnlineShops(begin, end int, where, order string) (int, []*dto.ListOnlineShop) {
	n, rows := si.query.PagedOnBusinessOnlineShops(begin, end, where, order)
	if len(rows) > 0 {
		for _, v := range rows {
			v.Logo = format.GetResUrl(v.Logo)
			if v.Host == "" {
				v.Host = v.Alias + "." + variable.Domain
			}
		}
	}
	return n, rows
}

func (si *shopServiceImpl) parseShop(sp *shop.OnlineShop) *shop_service.SShop {
	return &shop_service.SShop{
		ID:         int32(sp.Id),
		VendorId:   int32(sp.VendorId),
		ShopName:   sp.ShopName,
		Alias:      sp.Alias,
		Host:       sp.Host,
		Logo:       sp.Logo,
		ShopTitle:  sp.ShopTitle,
		ShopNotice: sp.ShopNotice,
	}
}
