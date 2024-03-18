import {Button, Link, Modal, ModalBody, ModalContent, ModalHeader} from "@nextui-org/react";
import JobForm from "./jobForm.tsx";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faXmark} from "@fortawesome/free-solid-svg-icons";

type JobFormModalProps = {
    isOpen: boolean;
    onOpenChange: (isOpen: boolean) => void;
}

function JobFormModal ({isOpen, onOpenChange}: JobFormModalProps) {
    return (
        <Modal
            classNames={{
                base: "light text-foreground bg-background md:w-2/3 sm:w-full w-full right-0 top-0 h-full m-0 p-0",
                body: "p-3",
                header: "border-b-2 p-3",
            }}
            hideCloseButton={true}
            isOpen={isOpen}
            onOpenChange={onOpenChange}
            placement="top"
            size={"full"}
            motionProps={{
                variants: {
                    enter: {
                        x: 0,
                        opacity: [0, 1],
                        transition: {
                            duration: 0.3,
                            ease: "easeOut",
                        },
                    },
                    exit: {
                        x: "100%",
                        opacity: 0,
                        transition: {
                            duration: 0.2,
                            ease: "easeIn",
                        },
                    },
                }
            }}
        >
            <ModalContent className={"right-0 absolute"}>
                <ModalHeader>
                    <div className="h-full w-full flex flex-row items-center">
                        <div>
                            <Link onClick={() => onOpenChange(false)} className={"mr-2 m-auto"}>
                                <FontAwesomeIcon icon={faXmark} />
                            </Link>
                        </div>
                        <div className={"ml-3"}>
                            Add Job
                        </div>
                        <div className="ml-auto">
                            <Button variant="bordered" className="mr-3">
                                Cancel
                            </Button>
                            <Button color="primary" className={"text-white font-bold"}>
                                Create
                            </Button>
                        </div>
                    </div>
                </ModalHeader>
                <ModalBody>
                    <JobForm
                        apiErrors={
                            {
                                fields: {
                                    name: "Name is required",
                                    type: "Type is required",
                                    slackWebhook: "Slack webhook is required"
                                }
                            }
                        }
                        onAddJob={() => onOpenChange(false)} />
                </ModalBody>
            </ModalContent>
        </Modal>
    )
}

export default JobFormModal;
