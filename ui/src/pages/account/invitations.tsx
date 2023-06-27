import { useCallback } from 'react'
import {
  Center,
  HStack,
  IconButton,
  Menu,
  MenuButton,
  MenuItem,
  MenuList,
  Portal,
  Stack,
  Table,
  Tbody,
  Td,
  Text,
  Th,
  Thead,
  Tr,
  useToast,
} from '@chakra-ui/react'
import { IconDotsVertical, SectionSpinner, variables } from '@koupr/ui'
import { Helmet } from 'react-helmet-async'
import InvitationAPI, { Invitation, SortOrder } from '@/client/api/invitation'
import UserAPI from '@/client/idp/user'
import { swrConfig } from '@/client/options'
import PagePagination, {
  usePagePagination,
} from '@/components/common/page-pagination'

const AccountInvitationsPage = () => {
  const toast = useToast()
  const { data: user, error: userError } = UserAPI.useGet()
  const { page, size, onPageChange, onSizeChange } = usePagePagination({
    localStoragePrefix: 'voltaserve',
    localStorageNamespace: 'incoming_invitation',
  })
  const {
    data: list,
    error: invitationsError,
    mutate,
  } = InvitationAPI.useGetIncoming(
    { page, size, sortOrder: SortOrder.Desc },
    swrConfig()
  )

  const handleAccept = useCallback(
    async (invitationId: string) => {
      await InvitationAPI.accept(invitationId)
      mutate()
      toast({
        title: 'Invitation accepted',
        status: 'success',
        isClosable: true,
      })
    },
    [mutate, toast]
  )

  const handleDecline = useCallback(
    async (invitationId: string) => {
      await InvitationAPI.delete(invitationId)
      mutate()
      toast({
        title: 'Invitation declined',
        status: 'info',
        isClosable: true,
      })
    },
    [mutate, toast]
  )

  if (userError || invitationsError) {
    return null
  }
  if (!user || !list) {
    return <SectionSpinner />
  }

  return (
    <>
      <Helmet>
        <title>{user.fullName}</title>
      </Helmet>
      {list.data.length === 0 && (
        <Center h="300px">
          <Text>There are no invitations.</Text>
        </Center>
      )}
      {list.data.length > 0 && (
        <Stack
          direction="column"
          spacing={variables.spacing2Xl}
          pb={variables.spacing2Xl}
        >
          <Table variant="simple">
            <Thead>
              <Tr>
                <Th>From</Th>
                <Th>Organization</Th>
                <Th>Date</Th>
                <Th></Th>
              </Tr>
            </Thead>
            <Tbody>
              {list.data.length > 0 &&
                list.data.map((e: Invitation) => (
                  <Tr key={e.id}>
                    <Td>{e.owner.fullName}</Td>
                    <Td>{e.organization.name}</Td>
                    <Td>{e.updateTime || e.createTime}</Td>
                    <Td textAlign="right">
                      <Menu>
                        <MenuButton
                          as={IconButton}
                          icon={<IconDotsVertical />}
                          variant="ghost"
                          aria-label=""
                        />
                        <Portal>
                          <MenuList>
                            <MenuItem onClick={() => handleAccept(e.id)}>
                              Accept
                            </MenuItem>
                            <MenuItem
                              color="red"
                              onClick={() => handleDecline(e.id)}
                            >
                              Decline
                            </MenuItem>
                          </MenuList>
                        </Portal>
                      </Menu>
                    </Td>
                  </Tr>
                ))}
            </Tbody>
          </Table>
          {list && (
            <HStack alignSelf="end">
              <PagePagination
                totalPages={list.totalPages}
                page={page}
                size={size}
                onPageChange={onPageChange}
                onSizeChange={onSizeChange}
              />
            </HStack>
          )}
        </Stack>
      )}
    </>
  )
}

export default AccountInvitationsPage
