import React from "react";
import AddMediaPage from "../../src/pages/AddMediaPage";
import { useApolloClient } from "@apollo/client/react";
import { createMedia } from "../../lib/createMedia";
import { createActivity } from "../../lib/createActivity";
import { GET_MOVIES_QUERY, ME_ACTIVITIES_QUERY } from "../../lib/graphql";
import { Media } from "../../src/types";

export default function AddTabPage() {
  const apolloClient = useApolloClient();
  const [isAddingMedia, setIsAddingMedia] = React.useState(false);

  const handleAddMedia = async (newMedia: Omit<Media, "id">) => {
    setIsAddingMedia(true);
    try {
      const result = await createMedia(newMedia);
      if (result?.id) {
        await createActivity({
          mediaId: result.id,
          statusId: 1,
        });
      }
      await apolloClient.refetchQueries({ include: [GET_MOVIES_QUERY, ME_ACTIVITIES_QUERY] });
    } catch (error) {
      console.error("Failed to add media:", error);
    } finally {
      setIsAddingMedia(false);
    }
  };

  return <AddMediaPage onBack={() => {}} onAddMedia={handleAddMedia} isLoading={isAddingMedia} />;
}
