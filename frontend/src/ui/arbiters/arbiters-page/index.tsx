import React from 'react';
import TextCleaner from "../../common/participants-cleaner";
import EventRewardGiver from "../event-reward-giver";
import Protector from "../../protector";
import BusinessRewarder from "../business-rewarder";
import {PERMISSION} from "../../../model/user/types";
import {LayoutResponsiveGrid} from "../../../components/common/layouts/responsive-grid";
import {LayoutReponsiveFlex} from "../../../components/common/layouts/responsive-flex";
import ChecksTable from "../../economics/checks-table";
import ArbiterItemsTable from "../items-table";

const ArbitersPage = () => {
    return (
        <Protector accessLevel={PERMISSION.Arbiter}>
            <LayoutResponsiveGrid gap={10} columns={[{minWidth: 620}, {minWidth: 350}]}>
                <LayoutReponsiveFlex direction={"column"} gap={10}>
                    <TextCleaner/>
                    <BusinessRewarder/>
                    <ArbiterItemsTable/>
                </LayoutReponsiveFlex>
                <LayoutResponsiveGrid gap={10} columns={[{}]}>
                    <EventRewardGiver/>
                    <ChecksTable/>
                </LayoutResponsiveGrid>
            </LayoutResponsiveGrid>
        </Protector>
    );
};

export default React.memo(ArbitersPage);